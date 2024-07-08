package storage

import (
	"database/sql"
	"log/slog"
	"strings"

	_ "github.com/mattn/go-sqlite3"

	"github.com/tnikic/oni/model"
	"github.com/tnikic/oni/tools"
)

type SQLite struct {
	path string
	db   *sql.DB
}

func InitSQLite() *SQLite {
	config := tools.LoadConfig()

	sqlite := &SQLite{
		path: config.DB.Path,
	}

	db, err := sql.Open("sqlite3", sqlite.path)
	if err != nil {
		slog.Error("Error opening SQLite database")
		panic(err)
	}

	sqlite.db = db

	return sqlite
}

// --------------------
// Manga
// --------------------
func (s *SQLite) StoreManga(manga *model.Manga) {
	sqlStatement := `
        INSERT INTO manga (title, all_titles, author, artist, status, country, description, path)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING id
    `
	id := 0

	err := s.db.QueryRow(
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

func (s *SQLite) GetManga(mangaID int) *model.Manga {
	sqlStatement := `
        SELECT title, all_titles, author, artist, status, country, description, path
        FROM manga
        WHERE id = $1
    `
	row := s.db.QueryRow(sqlStatement, mangaID)

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

func (s *SQLite) DeleteManga(mangaID int) {
	sqlStatement := `
        DELETE FROM manga
        WHERE id = $1
    `

	_, err := s.db.Exec(sqlStatement, mangaID)
	if err != nil {
		slog.Error("Error deleting manga from PostgreSQL")
	}
}

func (s *SQLite) GetAllManga() []*model.Manga {
	sqlStatement := `
        SELECT id, title, all_titles, author, artist, status, country, description, path
        FROM manga
    `

	rows, err := s.db.Query(sqlStatement)
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

func (s *SQLite) StoreChapter(chapter *model.Chapter) {
	sqlStatement := `
        INSERT INTO chapter (manga_id, chapter_entry_id, number)
        VALUES ($1, $2, $3)
        RETURNING id
    `
	id := 0

	err := s.db.QueryRow(
		sqlStatement,
		chapter.MangaID,
		chapter.ChapterEntryID,
		chapter.Number,
	).Scan(&id)
	if err != nil {
		slog.Error("Error storing chapter in PostgreSQL")
	}

	chapter.ID = id
}

func (s *SQLite) GetChapter(chapterID int) *model.Chapter {
	sqlStatement := `
        SELECT manga_id, chapter_entry_id, number
        FROM chapter
        WHERE id = $1
    `

	row := s.db.QueryRow(sqlStatement, chapterID)

	chapter := &model.Chapter{}
	err := row.Scan(
		&chapter.MangaID,
		&chapter.ChapterEntryID,
		&chapter.Number,
	)
	if err != nil {
		slog.Error("Error fetching chapter from PostgreSQL")
		return nil
	}

	chapter.ID = chapterID
	return chapter
}

func (s *SQLite) DeleteChapter(chapterID int) {
	sqlStatement := `
        DELETE FROM chapter
        WHERE id = $1
    `

	_, err := s.db.Exec(sqlStatement, chapterID)
	if err != nil {
		slog.Error("Error deleting chapter from PostgreSQL")
	}
}

func (s *SQLite) GetMangaChapters(mangaID int) []*model.Chapter {
	sqlStatement := `
        SELECT id, chapter_entry_id, number
        FROM chapter
        WHERE manga_id = $1
    `

	rows, err := s.db.Query(sqlStatement, mangaID)
	if err != nil {
		slog.Error("Error fetching chapters from PostgreSQL")
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
			slog.Error("Error scanning chapter from PostgreSQL")
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

func (s *SQLite) StorePlatform(platform *model.Platform) {
	sqlStatement := `
        INSERT INTO platform (name, active, ranking)
        VALUES ($1, $2, $3)
        RETURNING id
    `
	var id string

	err := s.db.QueryRow(
		sqlStatement,
		platform.Name,
		platform.Active,
		platform.Ranking,
	).Scan(&id)
	if err != nil {
		slog.Error("Error storing platform in PostgreSQL")
	}

	platform.ID = id
}

func (s *SQLite) GetPlatform(platformID string) *model.Platform {
	sqlStatement := `
        SELECT name, active, ranking
        FROM platform
        WHERE id = $1
    `

	row := s.db.QueryRow(sqlStatement, platformID)

	platform := &model.Platform{}
	err := row.Scan(
		&platform.Name,
		&platform.Active,
		&platform.Ranking,
	)
	if err != nil {
		slog.Error("Error fetching platform from PostgreSQL")
		return nil
	}

	platform.ID = platformID
	return platform
}

func (s *SQLite) DeletePlatform(platformID string) {
	sqlStatement := `
        DELETE FROM platform
        WHERE id = $1
    `

	_, err := s.db.Exec(sqlStatement, platformID)
	if err != nil {
		slog.Error("Error deleting platform from PostgreSQL")
	}
}

func (s *SQLite) GetAllPlatforms() []*model.Platform {
	sqlStatement := `
        SELECT id, name, active, ranking
        FROM platform
    `

	rows, err := s.db.Query(sqlStatement)
	if err != nil {
		slog.Error("Error fetching all platforms from PostgreSQL")
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
			slog.Error("Error scanning platform from PostgreSQL")
			return nil
		}

		platforms = append(platforms, platform)
	}

	return platforms
}

// --------------------
// Manga Entries
// --------------------

func (s *SQLite) StoreMangaEntry(mangaEntry *model.MangaEntry) {
	sqlStatement := `
        INSERT INTO manga_entry (platform_id, manga_id, name, manga_url, chapter_count, ranking)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id
    `
	id := 0

	err := s.db.QueryRow(
		sqlStatement,
		mangaEntry.PlatformID,
		mangaEntry.MangaID,
		mangaEntry.Name,
		mangaEntry.MangaUrl,
		mangaEntry.ChapterCount,
		mangaEntry.Ranking,
	).Scan(&id)
	if err != nil {
		slog.Error("Error storing manga entry in PostgreSQL")
	}

	mangaEntry.ID = id
}

func (s *SQLite) GetMangaEntry(mangaEntryID int) *model.MangaEntry {
	sqlStatement := `
        SELECT platform_id, manga_id, name, manga_url, chapter_count, ranking
        FROM manga_entry
        WHERE id = $1
    `

	row := s.db.QueryRow(sqlStatement, mangaEntryID)

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
		slog.Error("Error fetching manga entry from PostgreSQL")
		return nil
	}

	mangaEntry.ID = mangaEntryID
	return mangaEntry
}

func (s *SQLite) DeleteMangaEntry(mangaEntryID int) {
	sqlStatement := `
        DELETE FROM manga_entry
        WHERE id = $1
    `

	_, err := s.db.Exec(sqlStatement, mangaEntryID)
	if err != nil {
		slog.Error("Error deleting manga entry from PostgreSQL")
	}
}

func (s *SQLite) GetMangaEntries(mangaID int) []*model.MangaEntry {
	sqlStatement := `
        SELECT id, platform_id, name, manga_url, chapter_count, ranking
        FROM manga_entry
        WHERE manga_id = $1
    `

	rows, err := s.db.Query(sqlStatement, mangaID)
	if err != nil {
		slog.Error("Error fetching manga entries from PostgreSQL")
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
			slog.Error("Error scanning manga entry from PostgreSQL")
			return nil
		}

		entry.MangaID = mangaID
		entries = append(entries, entry)
	}

	return entries
}

func (s *SQLite) GetPlatformEntries(platformID string) []*model.MangaEntry {
	sqlStatement := `
        SELECT id, manga_id, name, manga_url, chapter_count, ranking
        FROM manga_entry
        WHERE platform_id = $1
    `

	rows, err := s.db.Query(sqlStatement, platformID)
	if err != nil {
		slog.Error("Error fetching manga entries from PostgreSQL")
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
			slog.Error("Error scanning manga entry from PostgreSQL")
			return nil
		}

		entry.PlatformID = platformID
		entries = append(entries, entry)
	}

	return entries
}

func (s *SQLite) GetAllMangaEntries() []*model.MangaEntry {
	sqlStatement := `
        SELECT id, platform_id, manga_id, name, manga_url, chapter_count, ranking
        FROM manga_entry
    `

	rows, err := s.db.Query(sqlStatement)
	if err != nil {
		slog.Error("Error fetching all manga entries from PostgreSQL")
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
			slog.Error("Error scanning manga entry from PostgreSQL")
			return nil
		}

		entries = append(entries, entry)
	}

	return entries
}

// --------------------
// Chapter Entries
// --------------------

func (s *SQLite) StoreChapterEntry(chapterEntry *model.ChapterEntry) {
	sqlStatement := `
        INSERT INTO chapter_entry (chapter_id, manga_entry_id, manga_id, chapter_number, chapter_url)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id
    `
	id := 0

	err := s.db.QueryRow(
		sqlStatement,
		chapterEntry.ChapterID,
		chapterEntry.MangaEntryID,
		chapterEntry.MangaID,
		chapterEntry.ChapterNumber,
		chapterEntry.ChapterUrl,
	).Scan(&id)
	if err != nil {
		slog.Error("Error storing chapter entry in PostgreSQL")
	}

	chapterEntry.ID = id
}

func (s *SQLite) GetChapterEntry(chapterEntryID int) *model.ChapterEntry {
	sqlStatement := `
        SELECT chapter_id, manga_entry_id, manga_id, chapter_number, chapter_url
        FROM chapter_entry
        WHERE id = $1
    `

	row := s.db.QueryRow(sqlStatement, chapterEntryID)

	chapterEntry := &model.ChapterEntry{}
	err := row.Scan(
		&chapterEntry.ChapterID,
		&chapterEntry.MangaEntryID,
		&chapterEntry.MangaID,
		&chapterEntry.ChapterNumber,
		&chapterEntry.ChapterUrl,
	)
	if err != nil {
		slog.Error("Error fetching chapter entry from PostgreSQL")
		return nil
	}

	chapterEntry.ID = chapterEntryID
	return chapterEntry
}

func (s *SQLite) DeleteChapterEntry(chapterEntryID int) {
	sqlStatement := `
        DELETE FROM chapter_entry
        WHERE id = $1
    `

	_, err := s.db.Exec(sqlStatement, chapterEntryID)
	if err != nil {
		slog.Error("Error deleting chapter entry from PostgreSQL")
	}
}

func (s *SQLite) GetChapterEntryForChapter(mangaEntryID int, chapterNumber int) *model.ChapterEntry {
	sqlStatement := `
        SELECT id, chapter_id, manga_entry_id, manga_id, chapter_number, chapter_url
        FROM chapter_entry
        WHERE manga_entry_id = $1 AND chapter_number = $2
    `

	row := s.db.QueryRow(sqlStatement, mangaEntryID, chapterNumber)

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
		slog.Error("Error fetching chapter entry for chapter from PostgreSQL")
		return nil
	}

	return chapterEntry
}

func (s *SQLite) GetChapterEntries(mangaEntryID int) []*model.ChapterEntry {
	sqlStatement := `
        SELECT id, chapter_id, manga_entry_id, manga_id, chapter_number, chapter_url
        FROM chapter_entry
        WHERE manga_entry_id = $1
    `

	rows, err := s.db.Query(sqlStatement, mangaEntryID)
	if err != nil {
		slog.Error("Error fetching chapter entries from PostgreSQL")
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
			slog.Error("Error scanning chapter entry from PostgreSQL")
			return nil
		}

		entries = append(entries, entry)
	}

	return entries
}

// --------------------
// Private Functions
// --------------------

func (s *SQLite) createTables() {
	_, err := s.db.Exec(`
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
        );
    `)
	if err != nil {
		slog.Error("Error creating manga table")
		panic(err)
	}

	_, err = s.db.Exec(`
        CREATE TABLE IF NOT EXISTS chapter (
            id SERIAL PRIMARY KEY,
            manga_id INT,
            chapter_entry_id INT,
            number INT,
        );
    `)
	if err != nil {
		slog.Error("Error creating chapter table")
		panic(err)
	}

	_, err = s.db.Exec(`
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

	_, err = s.db.Exec(`
        CREATE TABLE IF NOT EXISTS manga_entry (
            id SERIAL PRIMARY KEY,
            manga_id INT,
            platform_id INT,
            name TEXT,
            manga_url TEXT,
            chapter_count INT
            ranking INT,
        );
    `)
	if err != nil {
		slog.Error("Error creating manga entry table")
		panic(err)
	}

	_, err = s.db.Exec(`
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