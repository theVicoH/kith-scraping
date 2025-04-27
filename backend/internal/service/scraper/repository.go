package scraper

import (
	"database/sql"
	"fmt"
)

type ScrapeTarget struct {
	ID       int
	URL      string
	Category string
	Active   bool
}

func GetScrapeTargetByURL(db *sql.DB, url string) (*ScrapeTarget, error) {
	const q = `
	SELECT id, url, category, active
	FROM scrape_targets
	WHERE url = $1
	`
	var t ScrapeTarget
	err := db.QueryRow(q, url).
		Scan(&t.ID, &t.URL, &t.Category, &t.Active)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("GetScrapeTargetByURL: %w", err)
	}
	return &t, nil
}

func InsertScrapeTarget(db *sql.DB, url, category string) error {
	const q = `
	INSERT INTO scrape_targets (url, category, active)
	VALUES ($1, $2, true)
	`
	if _, err := db.Exec(q, url, category); err != nil {
		return fmt.Errorf("InsertScrapeTarget: %w", err)
	}
	return nil
}

func UpdateScrapeTarget(db *sql.DB, url, category string) error {
	const q = `
	UPDATE scrape_targets
	SET category = $1, active = true
	WHERE url = $2
	`
	if _, err := db.Exec(q, category, url); err != nil {
		return fmt.Errorf("UpdateScrapeTarget: %w", err)
	}
	return nil
}

func LoadActiveTargets(db *sql.DB) ([]ScrapeTarget, error) {
	const q = `
	SELECT id, url, category, active
	FROM scrape_targets
	WHERE active = true
	`
	rows, err := db.Query(q)
	if err != nil {
		return nil, fmt.Errorf("LoadActiveTargets: %w", err)
	}
	defer rows.Close()

	var list []ScrapeTarget
	for rows.Next() {
		var t ScrapeTarget
		if err := rows.Scan(&t.ID, &t.URL, &t.Category, &t.Active); err != nil {
			return nil, fmt.Errorf("LoadActiveTargets scan: %w", err)
		}
		list = append(list, t)
	}
	return list, nil
}
