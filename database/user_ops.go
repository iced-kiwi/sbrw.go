package database

import (
	"fmt"
	"gosbrw/database/structs"
	"log"
)

func CreateUserTable() error {
	db := GetDB()
	if db == nil {
		return fmt.Errorf("database connection is not initialized")
	}

	createTableQuery := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		remote_user_id INTEGER,
		security_token TEXT,
		-- user_id INTEGER,                 -- Removed as per request, 'id' will serve as the primary user identifier.
		email TEXT UNIQUE NOT NULL,      -- Emails should be unique and not null.
		password TEXT NOT NULL,          -- Store hashed passwords, never plaintext.
		date_created TEXT,               -- From structs.User.DateCreated. Consider TIMESTAMPTZ for actual date/time objects.
		banned BOOLEAN DEFAULT FALSE,
		premium BOOLEAN DEFAULT FALSE,
		admin BOOLEAN DEFAULT FALSE,
		game_launcher_certificate TEXT,
		game_launcher_hash TEXT,
		hidden_hwid TEXT,
		hwid TEXT,
		os_version TEXT,
		user_agent TEXT,
		ip_address TEXT,
		default_persona_idx INTEGER,
		active_persona_id INTEGER,       -- This would likely be a foreign key to a personas table if personas have their own IDs.
		locked BOOLEAN DEFAULT FALSE,
		selected_persona_index INTEGER,
		full_game_access BOOLEAN DEFAULT FALSE,
		complete BOOLEAN DEFAULT FALSE,
		last_auth_date TEXT,             -- From structs.User.LastAuthDate. Consider TIMESTAMPTZ.
		subscribe_msg BOOLEAN DEFAULT FALSE,
		address1 TEXT,
		address2 TEXT,
		country TEXT,
		dob TEXT,                        -- Date of Birth.
		email_status TEXT,
		first_name TEXT,
		gender TEXT,
		id_digits TEXT,
		landline_phone TEXT,
		language TEXT,
		last_name TEXT,
		mobile TEXT,
		nickname TEXT,
		postal_code TEXT,
		real_name TEXT,
		reason_code TEXT,
		starter_pack_entitlement_tag TEXT,
		status TEXT,
		tos_version TEXT,
		username TEXT UNIQUE,            -- Usernames are often unique.
		appear_offline BOOLEAN DEFAULT FALSE,
		decline_group_invitations INTEGER,
		decline_incoming_friend_requests BOOLEAN DEFAULT FALSE,
		decline_private_invite INTEGER,
		hide_offline_friends BOOLEAN DEFAULT FALSE,
		show_news_on_sign_in BOOLEAN DEFAULT FALSE,
		created_at TIMESTAMPTZ DEFAULT NOW(), -- Standard practice: timestamp of row creation.
		updated_at TIMESTAMPTZ DEFAULT NOW()  -- Standard practice: timestamp of last row update.
	);`

	_, err := db.Exec(createTableQuery)
	if err != nil {
		log.Printf("Error creating users table: %v", err)
		return fmt.Errorf("error creating users table: %w", err)
	}
	log.Println("Users table checked/created successfully.")
	return nil
}

func GetUserByEmail(email string) (structs.User, error) {
	db := GetDB()
	var user structs.User
	err := db.QueryRow("SELECT id, email, password FROM users WHERE email=$1", email).
		Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		return structs.User{}, err
	}
	return user, nil
}

func GetUserByID(userID int) (structs.User, error) {
	db := GetDB()
	var user structs.User
	err := db.QueryRow("SELECT id, email, password FROM users WHERE id=$1", userID).
		Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		return structs.User{}, err
	}
	return user, nil
}

func CreateNewUser(userInfo structs.User) (int, error) {
	db := GetDB()
	insertQuery := `
	INSERT INTO users (
		remote_user_id, security_token, email, password, date_created,
		banned, premium, admin, game_launcher_certificate, game_launcher_hash,
		hidden_hwid, hwid, os_version, user_agent, ip_address,
		default_persona_idx, active_persona_id, locked, selected_persona_index,
		full_game_access, complete, last_auth_date, subscribe_msg,
		address1, address2, country, dob, email_status, first_name, gender,
		id_digits, landline_phone, language, last_name, mobile, nickname,
		postal_code, real_name, reason_code, starter_pack_entitlement_tag,
		status, tos_version, username, appear_offline, decline_group_invitations,
		decline_incoming_friend_requests, decline_private_invite,
		hide_offline_friends, show_news_on_sign_in
	) VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17,
		$18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31,
		$32, $33, $34, $35, $36, $37, $38, $39, $40, $41, $42, $43, $44, $45,
		$46, $47, $48, $49
	) RETURNING id;`

	var newUserID int

	err := db.QueryRow(
		insertQuery,
		userInfo.RemoteUserID, userInfo.SecurityToken, userInfo.Email, userInfo.Password, userInfo.DateCreated,
		userInfo.Banned, userInfo.Premium, userInfo.Admin, userInfo.GameLauncherCertificate, userInfo.GameLauncherHash,
		userInfo.HiddenHWID, userInfo.HWID, userInfo.OSVersion, userInfo.UserAgent, userInfo.IPAddress,
		userInfo.DefaultPersonaIDX, userInfo.ActivePersonaID, userInfo.Locked, userInfo.SelectedPersonaIndex,
		userInfo.FullGameAccess, userInfo.Complete, userInfo.LastAuthDate, userInfo.SubscribeMsg,
		userInfo.Address1, userInfo.Address2, userInfo.Country, userInfo.DOB, userInfo.EmailStatus, userInfo.FirstName, userInfo.Gender,
		userInfo.IDDigits, userInfo.LandlinePhone, userInfo.Language, userInfo.LastName, userInfo.Mobile, userInfo.Nickname,
		userInfo.PostalCode, userInfo.RealName, userInfo.ReasonCode, userInfo.StarterPackEntitlementTag,
		userInfo.Status, userInfo.TOSVersion, userInfo.Username, userInfo.AppearOffline, userInfo.DeclineGroupInvitations,
		userInfo.DeclineIncomingFriendRequests, userInfo.DeclinePrivateInvite,
		userInfo.HideOfflineFriends, userInfo.ShowNewsOnSignIn,
	).Scan(&newUserID)

	if err != nil {
		log.Printf("Error creating new user with email %s: %v", userInfo.Email, err)
		return 0, fmt.Errorf("error creating new user: %w", err)
	}

	log.Printf("Successfully created new user with ID: %d, Email: %s", newUserID, userInfo.Email)
	return newUserID, nil
}
