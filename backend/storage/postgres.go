package storage

import (
	"database/sql"
	"fmt"
	"log/slog"
	"strings"

	_ "github.com/lib/pq"

	"github.com/tnikic/oni/model"
	"github.com/tnikic/oni/tools"
)

type PostgreSQL struct {
	db     *sql.DB
	config *tools.Config
}

func InitPostgreSQL() *PostgreSQL {
	postgres := PostgreSQL{
		config: tools.LoadConfig(),
	}
	connection := fmt.Sprintf(
		"postgresql://%s:%s@%s/%s?sslmode=disable",
		postgres.config.DB.Postgres.User,
		postgres.config.DB.Postgres.Password,
		postgres.config.DB.Postgres.Host,
		postgres.config.DB.Postgres.Database,
	)

	db, err := sql.Open(
		"postgres",
		connection,
	)
	if err != nil || db.Ping() != nil {
		slog.Error("Could not connect to the Postgres database")
		return nil
	}

	postgres.db = db
	postgres.createTables()
	return &postgres
}

// --------------------
// Manga
// --------------------
func (p *PostgreSQL) StoreManga(manga *model.Manga) {
	sqlStatement := `
        INSERT INTO manga (title, all_titles, author, artist, status, country, description, path)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING id
    `
	id := 0

	err := p.db.QueryRow(
		sqlStatement,
		manga.Title,
		strings.Join(manga.AllTitles, ","),
		manga.Author,
		manga.Artist,
		manga.Status,
		manga.Country,
		manga.Description,
		manga.Path,
	).Scan(&id)
	if err != nil {
		slog.Error("Error storing manga in PostgreSQL")
	}

	manga.ID = id
}

func (p *PostgreSQL) GetManga(mangaID int) *model.Manga {
	sqlStatement := `
        SELECT title, all_titles, author, artist, status, country, description, path
        FROM manga
        WHERE id = $1
    `
	row := p.db.QueryRow(sqlStatement, mangaID)

	var allTitles string
	manga := &model.Manga{}

	err := row.Scan(
		&manga.Title,
		&allTitles,
		&manga.Author,
		&manga.Artist,
		&manga.Status,
		&manga.Country,
		&manga.Description,
		&manga.Path,
	)
	if err != nil {
		slog.Error("Error fetching manga from PostgreSQL")
	}

	manga.AllTitles = strings.Split(allTitles, ",")
	manga.ID = mangaID

	return manga
}

func (p *PostgreSQL) DeleteManga(mangaID int) {
	sqlStatement := `
        DELETE FROM manga
        WHERE id = $1
    `

	_, err := p.db.Exec(sqlStatement, mangaID)
	if err != nil {
		slog.Error("Error deleting manga from PostgreSQL")
	}
}

func (p *PostgreSQL) GetAllManga() []*model.Manga {
	sqlStatement := `
        SELECT id, title, all_titles, author, artist, status, country, description, path
        FROM manga
    `

	rows, err := p.db.Query(sqlStatement)
	if err != nil {
		slog.Error("Error fetching all manga from PostgreSQL")
		return nil
	}
	defer rows.Close()

	var allTitles string
	mangas := []*model.Manga{}

	for rows.Next() {
		manga := &model.Manga{}
		err := rows.Scan(
			&manga.ID,
			&manga.Title,
			&allTitles,
			&manga.Author,
			&manga.Artist,
			&manga.Status,
			&manga.Country,
			&manga.Description,
			&manga.Path,
		)
		if err != nil {
			slog.Error("Error scanning manga from PostgreSQL")
			return nil
		}

		manga.AllTitles = strings.Split(allTitles, ",")

		mangas = append(mangas, manga)
	}

	return mangas
}

// --------------------
// Chapters
// --------------------

func (p *PostgreSQL) StoreChapter(chapter *model.Chapter) {
	sqlStatement := `
        INSERT INTO chapter (manga_id, chapter_entry_id, number)
        VALUES ($1, $2, $3)
        RETURNING id
    `
	id := 0

	err := p.db.QueryRow(
		sqlStatement,
		chapter.MangaID,
		chapter.ChapterEntryID,
		chapter.Number,
	).Scan(&id)
	if err != nil {
		slog.Error("Error storing chapter in sqlite")
	}

	chapter.ID = id
}

func (p *PostgreSQL) GetChapter(chapterID int) *model.Chapter {
	sqlStatement := `
        SELECT manga_id, chapter_entry_id, number
        FROM chapter
        WHERE id = $1
    `

	row := p.db.QueryRow(sqlStatement, chapterID)

	chapter := &model.Chapter{}
	err := row.Scan(
		&chapter.MangaID,
		&chapter.ChapterEntryID,
		&chapter.Number,
	)
	if err != nil {
		slog.Error("Error fetching chapter from sqlite")
		return nil
	}

	chapter.ID = chapterID
	return chapter
}

func (p *PostgreSQL) DeleteChapter(chapterID int) {
	sqlStatement := `
        DELETE FROM chapter
        WHERE id = $1
    `

	_, err := p.db.Exec(sqlStatement, chapterID)
	if err != nil {
		slog.Error("Error deleting chapter from sqlite")
	}
}

func (p *PostgreSQL) GetMangaChapters(mangaID int) []*model.Chapter {
	sqlStatement := `
        SELECT id, chapter_entry_id, number
        FROM chapter
        WHERE manga_id = $1
    `

	rows, err := p.db.Query(sqlStatement, mangaID)
	if err != nil {
		slog.Error("Error fetching chapters from sqlite")
		return nil
	}
	defer rows.Close()

	chapters := []*model.Chapter{}
	for rows.Next() {
		chapter := &model.Chapter{}
		err := rows.Scan(
			&chapter.ID,
			&chapter.ChapterEntryID,
			&chapter.Number,
		)
		if err != nil {
			slog.Error("Error scanning chapter from sqlite")
			return nil
		}

		chapter.MangaID = mangaID
		chapters = append(chapters, chapter)
	}

	return chapters
}

// --------------------
// Platforms
// --------------------

func (p *PostgreSQL) StorePlatform(platform *model.Platform) {
	sqlStatement := `
        INSERT INTO platform (id, name, active, ranking)
        VALUES ($1, $2, $3, $4)
    `

	_, err := p.db.Exec(
		sqlStatement,
		platform.ID,
		platform.Name,
		platform.Active,
		platform.Ranking,
	)
	if err != nil {
		slog.Error("Error storing platform in sqlite")
	}
}

func (p *PostgreSQL) GetPlatform(platformID string) *model.Platform {
	sqlStatement := `
        SELECT name, active, ranking
        FROM platform
        WHERE id = $1
    `

	row := p.db.QueryRow(sqlStatement, platformID)

	platform := &model.Platform{}
	err := row.Scan(
		&platform.Name,
		&platform.Active,
		&platform.Ranking,
	)
	if err != nil {
		slog.Error("Error fetching platform from sqlite")
		return nil
	}

	platform.ID = platformID
	return platform
}

func (p *PostgreSQL) DeletePlatform(platformID string) {
	sqlStatement := `
        DELETE FROM platform
        WHERE id = $1
    `

	_, err := p.db.Exec(sqlStatement, platformID)
	if err != nil {
		slog.Error("Error deleting platform from sqlite")
	}
}

func (p *PostgreSQL) GetAllPlatforms() []*model.Platform {
	sqlStatement := `
        SELECT id, name, active, ranking
        FROM platform
    `

	rows, err := p.db.Query(sqlStatement)
	if err != nil {
		slog.Error("Error fetching all platforms from sqlite")
		return nil
	}
	defer rows.Close()

	platforms := []*model.Platform{}
	for rows.Next() {
		platform := &model.Platform{}
		err := rows.Scan(
			&platform.ID,
			&platform.Name,
			&platform.Active,
			&platform.Ranking,
		)
		if err != nil {
			slog.Error("Error scanning platform from sqlite")
			return nil
		}

		platforms = append(platforms, platform)
	}

	return platforms
}

// --------------------
// Manga Entries
// --------------------

func (p *PostgreSQL) StoreMangaEntry(mangaEntry *model.MangaEntry) {
	sqlStatement := `
        INSERT INTO manga_entry (platform_id, manga_id, name, manga_url, chapter_count, ranking)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id
    `
	var id int

	err := p.db.QueryRow(
		sqlStatement,
		mangaEntry.PlatformID,
		mangaEntry.MangaID,
		mangaEntry.Name,
		mangaEntry.MangaUrl,
		mangaEntry.ChapterCount,
		mangaEntry.Ranking,
	).Scan(&id)
	if err != nil {
		slog.Error("Error storing manga entry in sqlite")
	}

	mangaEntry.ID = id
}

func (p *PostgreSQL) GetMangaEntry(mangaEntryID int) *model.MangaEntry {
	sqlStatement := `
        SELECT platform_id, manga_id, name, manga_url, chapter_count, ranking
        FROM manga_entry
        WHERE id = $1
    `

	row := p.db.QueryRow(sqlStatement, mangaEntryID)

	mangaEntry := &model.MangaEntry{}
	err := row.Scan(
		&mangaEntry.PlatformID,
		&mangaEntry.MangaID,
		&mangaEntry.Name,
		&mangaEntry.MangaUrl,
		&mangaEntry.ChapterCount,
		&mangaEntry.Ranking,
	)
	if err != nil {
		slog.Error("Error fetching manga entry from sqlite")
		return nil
	}

	mangaEntry.ID = mangaEntryID
	return mangaEntry
}

func (p *PostgreSQL) DeleteMangaEntry(mangaEntryID int) {
	sqlStatement := `
        DELETE FROM manga_entry
        WHERE id = $1
    `

	_, err := p.db.Exec(sqlStatement, mangaEntryID)
	if err != nil {
		slog.Error("Error deleting manga entry from sqlite")
	}
}

func (p *PostgreSQL) GetMangaEntries(mangaID int) []*model.MangaEntry {
	sqlStatement := `
        SELECT id, platform_id, name, manga_url, chapter_count, ranking
        FROM manga_entry
        WHERE manga_id = $1
    `

	rows, err := p.db.Query(sqlStatement, mangaID)
	if err != nil {
		slog.Error("Error fetching manga entries from sqlite")
	}
	defer rows.Close()

	entries := []*model.MangaEntry{}

	for rows.Next() {
		entry := &model.MangaEntry{}
		err := rows.Scan(
			&entry.ID,
			&entry.PlatformID,
			&entry.Name,
			&entry.MangaUrl,
			&entry.ChapterCount,
			&entry.Ranking,
		)
		if err != nil {
			slog.Error("Error scanning manga entry from sqlite")
			return nil
		}

		entry.MangaID = mangaID
		entries = append(entries, entry)
	}

	return entries
}

func (p *PostgreSQL) GetPlatformEntries(platformID string) []*model.MangaEntry {
	sqlStatement := `
        SELECT id, manga_id, name, manga_url, chapter_count, ranking
        FROM manga_entry
        WHERE platform_id = $1
    `

	rows, err := p.db.Query(sqlStatement, platformID)
	if err != nil {
		slog.Error("Error fetching manga entries from sqlite")
		return nil
	}
	defer rows.Close()

	entries := []*model.MangaEntry{}
	for rows.Next() {
		entry := &model.MangaEntry{}
		err := rows.Scan(
			&entry.ID,
			&entry.MangaID,
			&entry.Name,
			&entry.MangaUrl,
			&entry.ChapterCount,
			&entry.Ranking,
		)
		if err != nil {
			slog.Error("Error scanning manga entry from sqlite")
			return nil
		}

		entry.PlatformID = platformID
		entries = append(entries, entry)
	}

	return entries
}

func (p *PostgreSQL) GetAllMangaEntries() []*model.MangaEntry {
	sqlStatement := `
        SELECT id, platform_id, manga_id, name, manga_url, chapter_count, ranking
        FROM manga_entry
    `

	rows, err := p.db.Query(sqlStatement)
	if err != nil {
		slog.Error("Error fetching all manga entries from sqlite")
		return nil
	}
	defer rows.Close()

	entries := []*model.MangaEntry{}
	for rows.Next() {
		entry := &model.MangaEntry{}
		err := rows.Scan(
			&entry.ID,
			&entry.PlatformID,
			&entry.MangaID,
			&entry.Name,
			&entry.MangaUrl,
			&entry.ChapterCount,
			&entry.Ranking,
		)
		if err != nil {
			slog.Error("Error scanning manga entry from sqlite")
			return nil
		}

		entries = append(entries, entry)
	}

	return entries
}

// --------------------
// Chapter Entries
// --------------------

func (p *PostgreSQL) StoreChapterEntry(chapterEntry *model.ChapterEntry) {
	sqlStatement := `
        INSERT INTO chapter_entry (chapter_id, manga_entry_id, manga_id, chapter_number, chapter_url)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id
    `
	id := 0

	err := p.db.QueryRow(
		sqlStatement,
		chapterEntry.ChapterID,
		chapterEntry.MangaEntryID,
		chapterEntry.MangaID,
		chapterEntry.ChapterNumber,
		chapterEntry.ChapterUrl,
	).Scan(&id)
	if err != nil {
		slog.Error("Error storing chapter entry in sqlite")
	}

	chapterEntry.ID = id
}

func (p *PostgreSQL) GetChapterEntry(chapterEntryID int) *model.ChapterEntry {
	sqlStatement := `
        SELECT chapter_id, manga_entry_id, manga_id, chapter_number, chapter_url
        FROM chapter_entry
        WHERE id = $1
    `

	row := p.db.QueryRow(sqlStatement, chapterEntryID)

	chapterEntry := &model.ChapterEntry{}
	err := row.Scan(
		&chapterEntry.ChapterID,
		&chapterEntry.MangaEntryID,
		&chapterEntry.MangaID,
		&chapterEntry.ChapterNumber,
		&chapterEntry.ChapterUrl,
	)
	if err != nil {
		slog.Error("Error fetching chapter entry from sqlite")
		return nil
	}

	chapterEntry.ID = chapterEntryID
	return chapterEntry
}

func (p *PostgreSQL) DeleteChapterEntry(chapterEntryID int) {
	sqlStatement := `
        DELETE FROM chapter_entry
        WHERE id = $1
    `

	_, err := p.db.Exec(sqlStatement, chapterEntryID)
	if err != nil {
		slog.Error("Error deleting chapter entry from sqlite")
	}
}

func (p *PostgreSQL) GetChapterEntryForChapter(mangaEntryID int, chapterNumber int) *model.ChapterEntry {
	sqlStatement := `
        SELECT id, chapter_id, manga_entry_id, manga_id, chapter_number, chapter_url
        FROM chapter_entry
        WHERE manga_entry_id = $1 AND chapter_number = $2
    `

	row := p.db.QueryRow(sqlStatement, mangaEntryID, chapterNumber)

	chapterEntry := &model.ChapterEntry{}
	err := row.Scan(
		&chapterEntry.ID,
		&chapterEntry.ChapterID,
		&chapterEntry.MangaEntryID,
		&chapterEntry.MangaID,
		&chapterEntry.ChapterNumber,
		&chapterEntry.ChapterUrl,
	)
	if err != nil {
		slog.Error("Error fetching chapter entry for chapter from sqlite")
		return nil
	}

	return chapterEntry
}

func (p *PostgreSQL) GetChapterEntries(mangaEntryID int) []*model.ChapterEntry {
	sqlStatement := `
        SELECT id, chapter_id, manga_entry_id, manga_id, chapter_number, chapter_url
        FROM chapter_entry
        WHERE manga_entry_id = $1
    `

	rows, err := p.db.Query(sqlStatement, mangaEntryID)
	if err != nil {
		slog.Error("Error fetching chapter entries from sqlite")
		return nil
	}
	defer rows.Close()

	entries := []*model.ChapterEntry{}
	for rows.Next() {
		entry := &model.ChapterEntry{}
		err := rows.Scan(
			&entry.ID,
			&entry.ChapterID,
			&entry.MangaEntryID,
			&entry.MangaID,
			&entry.ChapterNumber,
			&entry.ChapterUrl,
		)
		if err != nil {
			slog.Error("Error scanning chapter entry from sqlite")
			return nil
		}

		entries = append(entries, entry)
	}

	return entries
}

// --------------------
// Private Functions
// --------------------

func (p *PostgreSQL) createTables() {
	_, err := p.db.Exec(`
        CREATE TABLE IF NOT EXISTS manga (
            id SERIAL PRIMARY KEY,
            title TEXT,
            all_titles TEXT,
            author TEXT,
            artist TEXT,
            status TEXT,
            country TEXT,
            description TEXT,
            path TEXT
        );`)
	if err != nil {
		slog.Error("Error creating manga table")
		panic(err)
	}

	_, err = p.db.Exec(`
        CREATE TABLE IF NOT EXISTS chapter (
            id SERIAL PRIMARY KEY,
            manga_id INT,
            chapter_entry_id INT,
            number INT
        );`)
	if err != nil {
		slog.Error("Error creating chapter table")
		panic(err)
	}

	_, err = p.db.Exec(`
        CREATE TABLE IF NOT EXISTS platform (
            id TEXT PRIMARY KEY,
            name TEXT,
            ranking INT,
            active BOOLEAN
        );
    `)
	if err != nil {
		slog.Error("Error creating platform table")
		panic(err)
	}

	_, err = p.db.Exec(`
        CREATE TABLE IF NOT EXISTS manga_entry (
            id SERIAL PRIMARY KEY,
            manga_id INT,
            platform_id TEXT,
            name TEXT,
            manga_url TEXT,
            chapter_count INT,
            ranking INT
        );
    `)
	if err != nil {
		slog.Error("Error creating manga entry table")
		panic(err)
	}

	_, err = p.db.Exec(`
        CREATE TABLE IF NOT EXISTS chapter_entry (
            id SERIAL PRIMARY KEY,
            chapter_id INT,
            manga_entry_id INT,
            manga_id INT,
            chapter_number INT,
            chapter_url TEXT
        );
    `)
	if err != nil {
		slog.Error("Error creating chapter entry table")
		panic(err)
	}
}
