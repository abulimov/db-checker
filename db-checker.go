// Copyright 2016 Alexander Bulimov. All rights reserved.
// Use of this source code is governed by a MIT license

// Utility to perform queries on PostgreSQL/MySQL database and warn on query result
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/abulimov/db-checker/utils"

	"github.com/abulimov/db-checker/base"

	"github.com/fractalcat/nagiosplugin"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

var version = "0.2.0"

// set up cli vars
var argDBType = flag.String("dbtype", "postgres", "Type of DB, can be 'postgres' or 'mysql'")
var argDBHost = flag.String("dbhost", "localhost", "Address of DB")
var argDBPort = flag.Int("dbport", 5432, "Port of Postgres DB")
var argDBName = flag.String("dbname", "postgres", "Name of DB")
var argDBUser = flag.String("dbuser", "postgres", "User to connect to DB")
var argDBPassword = flag.String("dbpassword", "", "Password to connect to DB (can be also set as PGPASSWORD or MYSQL_PWD environment variable)")
var argDBParams = flag.String("dbparams", "", "Optional params to pass in connection string, in param=value format, as a comma-separated list")
var argReport = flag.String("report", "", "Path for report file in JSON format")
var argDiff = flag.Bool("diff", false, "Check only diff between report and current state, rewrites old report")
var argCritical = flag.Bool("critical", false, "Consider any problem as CRITICAL (default is WARNING)")
var argChecksDir = flag.String("checks", "", "Path to directory with checks")
var argConcurrentChecks = flag.Int("concurrent-checks", 5, "Limit concurrent executions of checks")
var versionFlag = flag.Bool("version", false, "print db-checker version and exit")

func main() {
	check := nagiosplugin.NewCheck()
	// If we exit early or panic() we'll still output a result.
	defer check.Finish()

	flag.Parse()

	if *versionFlag {
		fmt.Printf("db-checker version %s\n", version)
		os.Exit(0)
	}

	if *argDBType != "mysql" && *argDBType != "postgres" {
		check.Unknownf("Not valid db type %s!\n, use 'postgres' or 'mysql'", *argDBType)
	}

	dbPassword := *argDBPassword
	if dbPassword == "" {
		switch *argDBType {
		case "postgres":
			if os.Getenv("PGPASSWORD") != "" {
				dbPassword = os.Getenv("PGPASSWORD")
			}
		case "mysql":
			if os.Getenv("MYSQL_PWD") != "" {
				dbPassword = os.Getenv("MYSQL_PWD")
			}
		}
	}

	if *argChecksDir == "" {
		check.Unknownf("'checks' option is required!")
	}

	// we cannot create diff without report file path
	if *argDiff && *argReport == "" {
		check.Unknownf("Diff check could only be performed when report is specified")
	}

	// choose what checks we should execute
	checks, err := base.GetChecks(*argChecksDir)
	if err != nil {
		check.Unknownf("%s", err)
	}
	// set up default report file
	reportFile := *argReport

	// format dbConnString
	var dbConnString string
	switch *argDBType {
	case "postgres":
		// postgres://username:password@address:port/dbname?param=value
		dbConnString = fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
			*argDBUser, dbPassword, *argDBHost, *argDBPort, *argDBName)
	case "mysql":
		// username:password@protocol(address:port)/dbname?param=value
		dbConnString = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
			*argDBUser, dbPassword, *argDBHost, *argDBPort, *argDBName)
	}

	if *argDBParams != "" {
		dbConnString += fmt.Sprintf("?%s", *argDBParams)
	}

	results, err := base.RunChecks(*argDBType, dbConnString, checks, *argConcurrentChecks)
	if err != nil {
		check.Unknownf("%s", err.Error())
	}

	filteredResults := results
	// if it is diff check
	if *argDiff {
		// try to read old report
		oldResults, err := utils.ReadReportFile(reportFile)
		if err != nil {
			base.Error.Printf("Failed to read report file %s: %s\n, assuming first check, all problems are new",
				reportFile, err)
		}
		// calculate diff
		filteredResults = utils.DiffResults(oldResults, results)
	}
	problemsCount, report := utils.ReportProblems(filteredResults)

	// Add some perfdata (label, unit, value, min, max, warn, crit).
	// The math.Inf(1) will be parsed as 'no maximum'.
	check.AddPerfDatum("problems", "", float64(problemsCount), 0.0, float64(0),
		float64(0), float64(0))

	// if we have no problems
	if problemsCount == 0 {
		check.AddResult(nagiosplugin.OK, report)
	} else {
		if *argCritical {
			check.AddResultf(nagiosplugin.CRITICAL, report)
		} else {
			check.AddResultf(nagiosplugin.WARNING, report)
		}

	}

	// write report
	if reportFile != "" {
		err := utils.WriteReportFile(results, reportFile)
		if err != nil {
			base.Error.Printf("Failed to generate report: %v\n", err)
		}
	}
}
