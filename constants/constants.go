package constants

const RequestId = "request_id"

const From = "X-Source-Slug"
const To = "X-Destination-Slug"
const RapidUrl = "X-Rapid-Url"
const KeyVersion = "X-Key-Version"

const Application = "application"
const ApplicationUlid = "application_ulid"

const RapidBridgeData = "./_rapid_bridge_data"

const EncryptionKeyValidityPeriod = 90 // in days
const SigningKeyValidityPeriod = 365   // in days

const RSAPrivateKeyFile = "rsa_private_key.pem"
const RSAPublicKeyFile = "rsa_public_key.pem"
const Ed25519PrivateKeyFile = "ed25519_private_key.pem"
const Ed25519PublicKeyFile = "ed25519_public_key.pem"
