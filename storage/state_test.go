package storage_test

import (
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pivotal-cf-experimental/bosh-bootloader/storage"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Store", func() {
	var (
		store   storage.Store
		tempDir string
	)

	BeforeEach(func() {
		store = storage.NewStore()

		var err error
		tempDir, err = ioutil.TempDir("", "")
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		storage.ResetEncode()
	})

	Describe("Set", func() {
		It("stores the state into a file", func() {
			err := store.Set(tempDir, storage.State{
				AWS: storage.AWS{
					AccessKeyID:     "some-aws-access-key-id",
					SecretAccessKey: "some-aws-secret-access-key",
					Region:          "some-region",
				},
				KeyPair: &storage.KeyPair{
					Name:       "some-name",
					PrivateKey: "some-private",
					PublicKey:  "some-public",
				},
				BOSH: &storage.BOSH{
					DirectorUsername:       "some-director-username",
					DirectorPassword:       "some-director-password",
					DirectorSSLCertificate: "some-bosh-ssl-certificate",
					DirectorSSLPrivateKey:  "some-bosh-ssl-private-key",
					State: map[string]interface{}{
						"key": "value",
					},
					Credentials: map[string]string{
						"mbusUsername":              "some-mbus-username",
						"natsUsername":              "some-nats-username",
						"postgresUsername":          "some-postgres-username",
						"registryUsername":          "some-registry-username",
						"blobstoreDirectorUsername": "some-blobstore-director-username",
						"blobstoreAgentUsername":    "some-blobstore-agent-username",
						"hmUsername":                "some-hm-username",
						"mbusPassword":              "some-mbus-password",
						"natsPassword":              "some-nats-password",
						"redisPassword":             "some-redis-password",
						"postgresPassword":          "some-postgres-password",
						"registryPassword":          "some-registry-password",
						"blobstoreDirectorPassword": "some-blobstore-director-password",
						"blobstoreAgentPassword":    "some-blobstore-agent-password",
						"hmPassword":                "some-hm-password",
					},
				},
				Stack: storage.Stack{
					Name: "some-stack-name",
				},
			})
			Expect(err).NotTo(HaveOccurred())

			data, err := ioutil.ReadFile(filepath.Join(tempDir, "state.json"))
			Expect(err).NotTo(HaveOccurred())
			Expect(data).To(MatchJSON(`{
				"version": 1,
				"aws": {
					"accessKeyId": "some-aws-access-key-id",
					"secretAccessKey": "some-aws-secret-access-key",
					"region": "some-region"
				},
				"keyPair": {
					"name": "some-name",
					"privateKey": "some-private",
					"publicKey": "some-public"
				},
				"bosh":{
					"directorUsername": "some-director-username",
					"directorPassword": "some-director-password",
					"directorSSLCertificate": "some-bosh-ssl-certificate",
					"directorSSLPrivateKey": "some-bosh-ssl-private-key",
					"credentials": {
						"mbusUsername": "some-mbus-username",
						"natsUsername": "some-nats-username",
						"postgresUsername": "some-postgres-username",
						"registryUsername": "some-registry-username",
						"blobstoreDirectorUsername": "some-blobstore-director-username",
						"blobstoreAgentUsername": "some-blobstore-agent-username",
						"hmUsername": "some-hm-username",
						"mbusPassword": "some-mbus-password",
						"natsPassword": "some-nats-password",
						"redisPassword": "some-redis-password",
						"postgresPassword": "some-postgres-password",
						"registryPassword": "some-registry-password",
						"blobstoreDirectorPassword": "some-blobstore-director-password",
						"blobstoreAgentPassword": "some-blobstore-agent-password",
						"hmPassword": "some-hm-password"
					},
					"state": {
						"key": "value"
					}
				},
				"stack": {
					"name": "some-stack-name"
				}
			}`))

			fileInfo, err := os.Stat(filepath.Join(tempDir, "state.json"))
			Expect(err).NotTo(HaveOccurred())
			Expect(fileInfo.Mode()).To(Equal(os.FileMode(0644)))
		})

		Context("failure cases", func() {
			It("fails to open the state.json file", func() {
				err := os.Chmod(tempDir, 0000)
				Expect(err).NotTo(HaveOccurred())

				err = store.Set(tempDir, storage.State{})
				Expect(err).To(MatchError(ContainSubstring("permission denied")))
			})

			It("fails to write the state.json file", func() {
				storage.SetEncode(func(io.Writer, interface{}) error {
					return errors.New("failed to encode")
				})

				err := store.Set(tempDir, storage.State{})
				Expect(err).To(MatchError("failed to encode"))
			})
		})
	})

	Describe("Get", func() {
		It("returns the stored state information", func() {
			err := ioutil.WriteFile(filepath.Join(tempDir, "state.json"), []byte(`{
				"version": 1,
				"aws": {
					"accessKeyId": "some-aws-access-key-id",
					"secretAccessKey": "some-aws-secret-access-key",
					"region": "some-aws-region"
				},
				"keyPair": {
					"name": "some-name",
					"privateKey": "some-private-key",
					"publicKey": "some-public-key"
				},
				"bosh": {
					"directorSSLCertificate": "some-bosh-ssl-certificate",
					"directorSSLPrivateKey": "some-bosh-ssl-private-key"
				},
				"stack": {
					"name": "some-stack-name"
				}
			}`), os.ModePerm)
			Expect(err).NotTo(HaveOccurred())

			state, err := store.Get(tempDir)
			Expect(err).NotTo(HaveOccurred())

			Expect(state).To(Equal(storage.State{
				Version: 1,
				AWS: storage.AWS{
					AccessKeyID:     "some-aws-access-key-id",
					SecretAccessKey: "some-aws-secret-access-key",
					Region:          "some-aws-region",
				},
				KeyPair: &storage.KeyPair{
					Name:       "some-name",
					PrivateKey: "some-private-key",
					PublicKey:  "some-public-key",
				},
				BOSH: &storage.BOSH{
					DirectorSSLCertificate: "some-bosh-ssl-certificate",
					DirectorSSLPrivateKey:  "some-bosh-ssl-private-key",
				},
				Stack: storage.Stack{
					Name: "some-stack-name",
				},
			}))
		})

		Context("when the state.json file doesn't exist", func() {
			It("returns an empty state object", func() {
				state, err := store.Get(tempDir)
				Expect(err).NotTo(HaveOccurred())

				Expect(state).To(Equal(storage.State{}))
			})
		})

		Context("failure cases", func() {
			It("fails to open the state.json file", func() {
				err := os.Chmod(tempDir, 0000)
				Expect(err).NotTo(HaveOccurred())

				_, err = store.Get(tempDir)
				Expect(err).To(MatchError(ContainSubstring("permission denied")))
			})

			It("fails to decode the state.json file", func() {
				err := ioutil.WriteFile(filepath.Join(tempDir, "state.json"), []byte(`%%%%`), os.ModePerm)
				Expect(err).NotTo(HaveOccurred())

				_, err = store.Get(tempDir)
				Expect(err).To(MatchError(ContainSubstring("invalid character")))
			})
		})
	})
})
