package util

const DefaultDRPolicy = "dr-policy"
const DefaultPlacement = "placement"

// Assume the first element in this array is the hub cluster
var ClusterNames = [3]string{"rdr-hub", "rdr-dr1", "rdr-dr2"}
