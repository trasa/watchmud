package db

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
	"github.com/trasa/watchmud/serverconfig"
	"golang.org/x/crypto/ssh"
	. "golang.org/x/crypto/ssh/agent"
	"io/ioutil"
	"net"
	"os"
	"time"
)

var watchdb *sqlx.DB

type ViaSSHDialer struct {
	client *ssh.Client
}

func (self *ViaSSHDialer) Open(s string) (_ driver.Conn, err error) {
	return pq.DialOpen(self, s)
}

func (self *ViaSSHDialer) Dial(network, address string) (net.Conn, error) {
	return self.client.Dial(network, address)
}

func (self *ViaSSHDialer) DialTimeout(network, address string, timeout time.Duration) (net.Conn, error) {
	return self.client.Dial(network, address)
}

func Init(config *serverconfig.Config) error {
	if config.DB.UseSSH {
		if err := initDBOverSSH(config); err != nil {
			return err
		}
	} else {
		if err := initDBDirectly(config); err != nil {
			return err
		}
	}
	if err := testConnection(); err != nil {
		return err
	}
	return nil
}

func initDBDirectly(config *serverconfig.Config) error {

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		config.DB.User,
		config.DB.Password,
		config.DB.Host,
		config.DB.Port,
		config.DB.Name)
	var err error
	watchdb, err = sqlx.Open("postgres", connStr)
	if err != nil {
		return err
	}
	return nil
}

func initDBOverSSH(config *serverconfig.Config) error {
	sshConfig, err := buildSSHConfig(config)
	if err != nil {
		log.Error().
			Err(err).
			Msg("Failed to construct SSH config")
		return err
	}
	// Connect to the SSH Server
	var sshcon *ssh.Client
	if sshcon, err = ssh.Dial("tcp", fmt.Sprintf("%s:%d", config.DB.SSH.Host, config.DB.SSH.Port), sshConfig); err != nil {
		log.Error().
			Err(err).
			Msgf("Failed to connect to SSH host %s:%d", config.DB.SSH.Host, config.DB.SSH.Port)
		return err
	}
	// sshcon is now connected
	//defer sshcon.Close()

	// Now we register the ViaSSHDialer with the ssh connection as a parameter
	sql.Register("postgres+ssh", &ViaSSHDialer{sshcon})

	// And now we can use our new driver with the regular postgres connection string tunneled through the SSH connection
	if db, err := sql.Open("postgres+ssh", fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", config.DB.User, config.DB.Password, config.DB.Host, config.DB.Port, config.DB.Name)); err == nil {
		log.Info().Msg("Successfully connected to DB over SSH")
		watchdb = sqlx.NewDb(db, "postgres")
		if err := watchdb.Ping(); err != nil {
			log.Error().Err(err).Msg("PING to postgres failed!")
			return err
		}
	}
	return nil
}

func testConnection() error {
	rows := watchdb.QueryRow("select now()")
	var now string
	if err := rows.Scan(&now); err != nil {
		return err
	}
	log.Info().Msgf("Database is live: %s", now)
	return nil
}

func buildSSHConfig(config *serverconfig.Config) (*ssh.ClientConfig, error) {
	var agentClient Agent
	// Establish a connection to the local ssh-agent
	if conn, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK")); err == nil {
		//defer conn.Close()

		// Create a new instance of the ssh agent
		agentClient = NewClient(conn)
	}
	// The client configuration with configuration option to use the ssh-agent
	sshConfig := &ssh.ClientConfig{
		User: config.DB.SSH.User,
		Auth: []ssh.AuthMethod{},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}
	// When the agentClient connection succeeded, add them as AuthMethod
	if agentClient != nil {
		sshConfig.Auth = append(sshConfig.Auth, ssh.PublicKeysCallback(agentClient.Signers))
	}

	keyAuth, err := getKeyAuthMethod(config.DB.SSH.KeyFile)
	if err != nil {
		return nil, err
	}
	sshConfig.Auth = append(sshConfig.Auth, keyAuth)
	return sshConfig, nil
}

func getKeyAuthMethod(keyfile string) (ssh.AuthMethod, error) {
	// add in the keyfile
	keybuf, err := ioutil.ReadFile(keyfile)
	if err != nil {
		return nil, err
	}
	key, err := ssh.ParsePrivateKey(keybuf)
	if err != nil {
		return nil, err
	}
	return ssh.PublicKeys(key), nil
}
