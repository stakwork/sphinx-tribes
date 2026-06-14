package db

import "time"

func (db database) ArchiveInactiveTribes(cutoffUnix int64) (int64, error) {
	if cutoffUnix <= 0 {
		return 0, nil
	}

	now := time.Now()
	result := db.db.Model(&Tribe{}).
		Where("(unlisted = ? OR unlisted IS NULL) AND (deleted = ? OR deleted IS NULL)", false, false).
		Where("last_active > 0 AND last_active < ?", cutoffUnix).
		Updates(map[string]interface{}{
			"unlisted": true,
			"updated":  &now,
		})

	return result.RowsAffected, result.Error
}
