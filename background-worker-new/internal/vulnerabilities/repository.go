package vulnerabilities

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/oss-health/background-worker/internal/db"
)

type VulnerabilityRepository struct {
	db *pgxpool.Pool
}

func NewVulnerabilityRepository(db *pgxpool.Pool) *VulnerabilityRepository {
	return &VulnerabilityRepository{db: db}
}

func (r *VulnerabilityRepository) InsertVulnerabilities(ctx context.Context, dependencyVersionVulnerabilities map[int][]Vulnerability) ([]db.Vulnerability, error) {
	if len(dependencyVersionVulnerabilities) == 0 {
		return nil, nil
	}

	queryArguments := []any{}
	clauses := []string{}
	paramIndex := 1

	for _, dependencyVersionVulnerability := range dependencyVersionVulnerabilities {
		for _, vuln := range dependencyVersionVulnerability {
			clauses = append(clauses, fmt.Sprintf("($%d, $%d, $%d)", paramIndex, paramIndex+1, paramIndex+2))
			queryArguments = append(queryArguments, vuln.ID, vuln.Summary, vuln.Severity)
			paramIndex += 3
		}
	}

	query := fmt.Sprintf(InsertVulnerabilityQuery, strings.Join(clauses, ", "))
	// logQueryWithParams(query, queryArguments)
	rows, err := r.db.Query(ctx, query, queryArguments...)
	if err != nil {
		return nil, fmt.Errorf("failed to insert dependencies: %v", err)
	}

	defer rows.Close()

	var inserted []db.Vulnerability
	for rows.Next() {
		var id int
		var osvID string
		if err := rows.Scan(&id, &osvID); err != nil {
			return nil, fmt.Errorf("failed to scan inserted row: %v", err)
		}
		inserted = append(inserted, struct {
			ID    int
			OSVID string
		}{ID: id, OSVID: osvID})
	}

	return inserted, nil
}

func (r *VulnerabilityRepository) InsertVersionVulnerability(ctx context.Context, versionVulnerabilities map[int]int) error {
	if len(versionVulnerabilities) == 0 {
		return nil
	}

	queryArguments := []any{}
	clauses := []string{}
	paramIndex := 1

	for vulnerabilityID, versionID := range versionVulnerabilities {
		clauses = append(clauses, fmt.Sprintf("($%d, $%d)", paramIndex, paramIndex+1))
		queryArguments = append(queryArguments, vulnerabilityID, versionID)
		paramIndex += 2
	}

	query := fmt.Sprintf(InsertVersionVulnerabilityQuery, strings.Join(clauses, ", "))
	// logQueryWithParams(query, queryArguments)
	_, err := r.db.Exec(ctx, query, queryArguments...)
	if err != nil {
		return fmt.Errorf("failed to insert version vulnerabilities: %v", err)
	}

	return nil
}


// func logQueryWithParams(query string, args []any) {
// 	for i, arg := range args {
// 		placeholder := fmt.Sprintf("$%d", i+1)
// 		var val string
// 		switch v := arg.(type) {
// 		case string:
// 			val = fmt.Sprintf("'%s'", strings.ReplaceAll(v, "'", "''")) // escape single quotes
// 		default:
// 			val = fmt.Sprintf("%v", v)
// 		}
// 		query = strings.Replace(query, placeholder, val, 1)
// 	}
// 	fmt.Println("Executing Query:", query)
// }
