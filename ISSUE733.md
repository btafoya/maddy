### [\#733] Dear maddy author, please tell me about the function of batch creating users.
**URL**: https://github.com/foxcpp/maddy/issues/733
**State**: CLOSED
**Labels**: new feature, 

**Body:**
# Use case

I currently have a requirement to have the golang program automatically create users. I tried to directly insert the database credentials.db, but the inserted database user could not log in normally. I would like to ask for guidance on how I should do this so that I can automatically generate users without having to manually go to the terminal to create new users.

![acd7b30297cc99c81313979e774b75dc](https://github.com/user-attachments/assets/a9b85b6b-760d-4183-8c4d-db16a67def8c)

I inserted this picture using the program, but I can't log in normally

---

## Resolution

**Root cause:**
Direct insertion into `credentials.db` without properly hashing the password and formatting the entry according to Maddy's `auth.pass_table` requirements will result in users being unable to log in. Maddy uses specific password hashing algorithms (by default, bcrypt) and a precise storage format (`hashAlgo:hashedPasswordString`) in the database. Simply inserting a plaintext password or an incorrectly hashed one will not work.

**Fix summary:**
Provided guidance on two methods for programmatic user creation, with a strong recommendation for using Maddy's CLI tool.

**How to reproduce (before/after if relevant):**
Before: Direct insertion of unhashed/incorrectly hashed passwords into `credentials.db` through a custom program leads to login failures.
After: Using the recommended `maddy creds add` command programmatically, or correctly implementing Maddy's hashing logic for direct database insertion, allows successful user creation and login.

**How to verify:**
1.  **Using `maddy creds add`:**
    ```bash
    maddy creds add testuser@example.com mysecretpassword
    ```
    Then try to authenticate a client with `testuser@example.com` and `mysecretpassword`.

2.  **Using direct DB insertion (if implemented):**
    Ensure the Maddy server is running and configured to use `auth.pass_table` with `sql_table`.
    Execute your Go program that performs the direct insertion with correct hashing.
    Then try to authenticate a client with the newly added user credentials.

**Notes / follow-ups:**

**Guidance for programmatic user creation:**

1.  **Recommended approach (using `maddy creds add` via `os/exec`):**
    The safest and most recommended way to programmatically add users is to use Maddy's own CLI tool, `maddy creds add`. This tool correctly handles password hashing (e.g., bcrypt with default cost) and updates the underlying `credentials.db` (or whatever `pass_table` is configured to use) in the format Maddy expects (`hashAlgo:hashedPasswordString`). Your Go program can execute this command using `os/exec`.

    **Example Go code snippet:**
    ```go
    package main

    import (
    	"fmt"
    	"os/exec"
    )

    func createUser(username, password string) error {
    	cmd := exec.Command("maddy", "creds", "add", username, password)
    	// Ensure maddy executable is in PATH or provide full path, e.g., "/usr/local/bin/maddy"
    	
    	output, err := cmd.CombinedOutput()
    	if err != nil {
    		return fmt.Errorf("failed to create user %s: %v\nOutput: %s", username, err, output)
    	}
    	fmt.Printf("User %s created successfully: %s\n", username, string(output))
    	return nil
    }

    func main() {
    	// Example usage
    	if err := createUser("newuser@example.com", "securepassword123"); err != nil {
    		fmt.Println("Error:", err)
    	}
    	if err := createUser("another@example.com", "anothersecret"); err != nil {
    		fmt.Println("Error:", err)
    	}
    }
    ```
    This method ensures compatibility and correctness with Maddy's internal user management logic.

2.  **Advanced approach (direct database insertion - use with caution):**
    If `os/exec` is not an option (e.g., due to sandboxing or security constraints), you can directly insert into `credentials.db` (assuming it's an SQLite database as is typical for `sql_table`). However, you *must* correctly hash the password using Maddy's algorithms and follow its storage format. The default hashing algorithm for `CreateUser` is `bcrypt`.

    *   **Password Hashing:** You will need to import `golang.org/x/crypto/bcrypt` into your Go program and use it to hash the password.
        ```go
        import (
        	"fmt"
        	"database/sql"
        	_ "github.com/mattn/go-sqlite3" // Import for SQLite driver
        	"golang.org/x/crypto/bcrypt"
        	"golang.org/x/text/secure/precis" // For username normalization
        )

        func insertUserDirect(dbPath, username, password string) error {
            db, err := sql.Open("sqlite3", dbPath)
            if err != nil {
                return fmt.Errorf("failed to open database: %w", err)
            }
            defer db.Close()

            // Normalize username using Maddy's logic
            key, err := precis.UsernameCaseMapped.CompareKey(username)
            if err != nil {
                return fmt.Errorf("failed to normalize username: %w", err)
            }

            // Hash password using bcrypt (default for Maddy)
            hashedPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
            if err != nil {
                return fmt.Errorf("failed to hash password: %w", err)
            }
            // Maddy stores hash as "bcrypt:hashedPasswordString"
            storedHash := "bcrypt:" + string(hashedPasswordBytes)

            // Insert into the 'passwords' table (common schema for sql_table)
            _, err = db.Exec("INSERT INTO passwords (username, password_hash) VALUES (?, ?)", key, storedHash)
            if err != nil {
                return fmt.Errorf("failed to insert user into DB: %w", err)
            }
            fmt.Printf("User %s inserted directly into %s\n", username, dbPath)
            return nil
        }

        func main() {
            // Example usage for direct DB insertion
            dbFilePath := "/var/lib/maddy/credentials.db" // Adjust to your Maddy credentials.db path
            if err := insertUserDirect(dbFilePath, "directuser@example.com", "directsecurepass"); err != nil {
                fmt.Println("Error:", err)
            }
        }
        ```
    *   **Username Normalization:** Maddy uses `golang.org/x/text/secure/precis.UsernameCaseMapped.CompareKey` to normalize usernames before storing them as keys. Your program *must* do the same to ensure the correct lookup.
    *   **Database Schema:** The `credentials.db` (SQLite) typically has a `passwords` table with at least two columns: `username` (or `key` from normalized username) and `password_hash`.

    This advanced approach requires deep knowledge of Maddy's internal mechanisms and is more fragile to Maddy's internal changes. The `os/exec` method is generally preferred.
