package database

import (
	"database/sql"
	"fmt"
	"gosbrw/database/structs"
	"log"

	"github.com/lib/pq"
)

func CreateServerInformationTable() error {
	db := GetDB()
	if db == nil {
		log.Fatal("Database connection is not initialized.")
		return fmt.Errorf("database connection is not initialized")
	}

	createTableQuery := `
	CREATE TABLE IF NOT EXISTS server_information (
		id SERIAL PRIMARY KEY,
		max_players INTEGER,
		online_players INTEGER,
		modern_auth_support BOOLEAN,
		message_srv TEXT,
		home_page_url TEXT,
		discord_url TEXT,
		facebook_url TEXT,
		twitter_url TEXT,
		server_name TEXT,
		country VARCHAR(255),
		timezone INTEGER,
		banner_url TEXT,
		admin_list TEXT[], 
		owner_list TEXT[], 
		number_of_registered INTEGER,
		seconds_to_shut_down INTEGER,
		activated_holiday_scenery_groups TEXT[],
		disactivated_holiday_scenery_groups TEXT[],
		require_ticket BOOLEAN,
		server_version VARCHAR(255),
		created_at TIMESTAMPTZ DEFAULT NOW(),
		updated_at TIMESTAMPTZ DEFAULT NOW()
	);`

	_, err := db.Exec(createTableQuery)
	if err != nil {
		log.Printf("Error creating server_information table: %v", err)
		return err
	}
	log.Println("server_information table checked/created successfully.")

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM server_information").Scan(&count)
	if err != nil {
		log.Printf("Error checking server_information table count: %v", err)
		return err
	}

	if count == 0 {
		log.Println("server_information table is empty. Populating with default data.")
		defaultInfo := structs.NewDefaultServerInformation()

		insertQuery := `
		INSERT INTO server_information (
			max_players, online_players, modern_auth_support, message_srv,
			home_page_url, discord_url, facebook_url, twitter_url,
			server_name, country, timezone, banner_url, admin_list, owner_list,
			number_of_registered, seconds_to_shut_down,
			activated_holiday_scenery_groups, disactivated_holiday_scenery_groups,
			require_ticket, server_version
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20
		)`

		_, err = db.Exec(insertQuery,
			defaultInfo.MaxPlayers, defaultInfo.OnlinePlayers, defaultInfo.ModernAuthSupport, defaultInfo.MessageSrv,
			defaultInfo.HomePageURL, defaultInfo.DiscordURL, defaultInfo.FacebookURL, defaultInfo.TwitterURL,
			defaultInfo.ServerName, defaultInfo.Country, defaultInfo.Timezone, defaultInfo.BannerURL,
			pq.Array(defaultInfo.AdminList), pq.Array(defaultInfo.OwnerList),
			defaultInfo.NumberOfRegistered, defaultInfo.SecondsToShutDown,
			pq.Array(defaultInfo.ActivatedHolidaySceneryGroups),
			pq.Array(defaultInfo.DisactivatedHolidaySceneryGroups),
			defaultInfo.RequireTicket, defaultInfo.ServerVersion,
		)
		if err != nil {
			log.Printf("Error inserting default server information: %v", err)
			return err
		}
		log.Println("Default server information inserted successfully.")
	} else {
		log.Println("server_information table already contains data.")
	}

	return nil
}

func GetServerInformation() (structs.ServerInformation, error) {
	db := GetDB()

	registeredCountQuery := "SELECT COUNT(*) FROM users"
	var registeredCount int
	err := db.QueryRow(registeredCountQuery).Scan(&registeredCount)
	if err != nil {
		log.Printf("Error getting registered count: %v", err)
		return structs.ServerInformation{}, err
	}
	updateQuery := "UPDATE server_information SET number_of_registered = $1"
	_, err = db.Exec(updateQuery, registeredCount)
	if err != nil {
		log.Printf("Error updating server_information table: %v", err)
		return structs.ServerInformation{}, err
	}
	
	var info structs.ServerInformation
	query := `
		SELECT
			max_players, online_players, modern_auth_support, message_srv,
			home_page_url, discord_url, facebook_url, twitter_url,
			server_name, country, timezone, banner_url,
			admin_list, owner_list, 
			number_of_registered, seconds_to_shut_down,
			activated_holiday_scenery_groups, disactivated_holiday_scenery_groups,
			require_ticket, server_version
		FROM server_information
		ORDER BY id ASC 
		LIMIT 1;`

	row := db.QueryRow(query)
	err = row.Scan(
		&info.MaxPlayers, &info.OnlinePlayers, &info.ModernAuthSupport, &info.MessageSrv,
		&info.HomePageURL, &info.DiscordURL, &info.FacebookURL, &info.TwitterURL,
		&info.ServerName, &info.Country, &info.Timezone, &info.BannerURL,
		pq.Array(&info.AdminList), pq.Array(&info.OwnerList),
		&info.NumberOfRegistered, &info.SecondsToShutDown,
		pq.Array(&info.ActivatedHolidaySceneryGroups), pq.Array(&info.DisactivatedHolidaySceneryGroups),
		&info.RequireTicket, &info.ServerVersion,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("No server information found in database.")
			return structs.ServerInformation{}, err
		}
		log.Printf("Error scanning server information: %v", err)
		return structs.ServerInformation{}, fmt.Errorf("error scanning server information: %w", err)
	}

	return info, nil
}
