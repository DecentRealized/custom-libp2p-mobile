package config

const MaxMessageSize = uint16(1 << 15) // Including Headers (16KB)

const MaxNotifierQueueSize = uint16(1 << 10) // Max Size = MaxMessageSize * MaxNotifierQueueSize (16MB)

const MdnsRendezvous = "file-drop-mdns-rendezvous"
