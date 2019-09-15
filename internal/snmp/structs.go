/*
Taken from https://github.com/wdreeveii/gopoller/blob/master/main.go to change and edit as necessary. This is for reference
and will most likely be completely rewritten for our methodology outlined in the scope.
 */
package snmp

import (
	"encoding/hex"
	"fmt"
	"github.com/coopernurse/gorp"
	"github.com/soniah/gosnmp"
	"math/rand"
	"time"
)

type SnmpPollingConfig struct {
	ResourceName                   string `json:"resource_name,omitempty" bson:"resource_name,omitempty"`
	Description                    string `json:"description,omitempty" bson:"description,omitempty"`
	IpAddress                      string `json:"ip_address,omitempty" bson:"ip_address,omitempty"`
	SnmpCommunityName              string `json:"snmp_community_name,omitempty" bson:"snmp_community_name,omitempty"`
	SnmpVersion                    string `json:"snmp_version,omitempty" bson:"snmp_version,omitempty"`
	SnmpV3SecurityLevel            string `json:"snmp_v_3_security_level,omitempty" bson:"snmp_v_3_security_level,omitempty"`
	SnmpV3AuthenticationProtocol   string `json:"snmp_v_3_authentication_protocol,omitempty" bson:"snmp_v_3_authentication_protocol,omitempty"`
	SnmpV3AuthenticationPassphrase string `json:"snmp_v_3_authentication_passphrase,omitempty" bson:"snmp_v_3_authentication_passphrase,omitempty"`
	SnmpV3PrivacyProtocol          string `json:"snmp_v_3_privacy_protocol,omitempty" bson:"snmp_v_3_privacy_protocol,omitempty"`
	SnmpV3PrivacyPassphrase        string `json:"snmp_v_3_privacy_passphrase,omitempty" bson:"snmp_v_3_privacy_passphrase,omitempty"`
	SnmpV3SecurityName             string `json:"snmp_v_3_security_name,omitempty" bson:"snmp_v_3_security_name,omitempty"`
	SnmpTimeout                    int    `json:"snmp_timeout,omitempty" bson:"snmp_timeout,omitempty"`
	SnmpRetries                    int    `json:"snmp_retries,omitempty" bson:"snmp_retries,omitempty"`
	SnmpEnabled                    string `json:"snmp_enabled,omitempty" bson:"snmp_enabled,omitempty"`
	Oid                            string `json:"oid,omitempty" bson:"oid,omitempty"`
	OidName                        string `json:"oid_name,omitempty" bson:"oid_name,omitempty"`
	PollType                       string `json:"poll_type,omitempty" bson:"poll_type,omitempty"`
	PollFreq                       int    `json:"poll_freq,omitempty" bson:"poll_freq,omitempty"`
	LastPollTime                   int64  `json:"last_poll_time,omitempty" bson:"last_poll_time,omitempty"`
	NextPollTime                   int64  `json:"next_poll_time,omitempty" bson:"next_poll_time,omitempty"`
	RealTimeReporting              string `json:"real_time_reporting,omitempty" bson:"real_time_reporting,omitempty"`
	History                        string `json:"history,omitempty" bson:"history,omitempty"`
}

func stringifyType(t gosnmp.Asn1BER) string {
	switch t {
	case gosnmp.Boolean:
		return "BOOLEAN"
	case gosnmp.Integer:
		return "INTEGER"
	case gosnmp.BitString:
		return "BITSTRING"
	case gosnmp.OctetString:
		return "OCTETSTRING"
	case gosnmp.Null:
		return "NULL"
	case gosnmp.ObjectIdentifier:
		return "OBJECTIDENTIFIER"
	case gosnmp.ObjectDescription:
		return "OBJECTDESCRIPTION"
	case gosnmp.IPAddress:
		return "IPADDRESS"
	case gosnmp.Counter32:
		return "COUNTER"
	case gosnmp.Gauge32:
		return "GAUGE"
	case gosnmp.TimeTicks:
		return "TIMETICKS"
	case gosnmp.Opaque:
		return "OPAQUE"
	case gosnmp.NsapAddress:
		return "NSAPADDRESS"
	case gosnmp.Counter64:
		return "COUNTER"
	case gosnmp.Uinteger32:
		return "UINTEGER"
	}
	return "UNKNOWN ASN1BER"
}


type SnmpFetchResult struct {
	Config SnmpPollingConfig
	Data   []gosnmp.SnmpPDU
	Err    error
}


// generate a bulk insert statement to insert the values
// into the database
func generateInsertData(res SnmpFetchResult) string {
	var data = make([]byte, 0, 50*len(res.Data))

	for i, v := range res.Data {
		if i != 0 {
			data = append(data, ", "...)
		}
		// convert byte arrays to hex encoded strings
		var value interface{}
		if nval, ok := v.Value.([]byte); ok {
			value = hex.EncodeToString(nval)
		} else {
			value = v.Value
		}
		data = append(data, "("...)
		data = append(data, fmt.Sprint(res.Config.LastPollTime/1000)...)
		data = append(data, ",'"...)
		data = append(data, res.Config.IpAddress...)
		data = append(data, "','"...)
		data = append(data, v.Name[1:]...)
		data = append(data, "','"...)
		data = append(data, stringifyType(v.Type)...)
		data = append(data, "','"...)
		data = append(data, fmt.Sprint(value)...)
		data = append(data, "')"...)
	}
	return string(data)
}


// do one snmp query
func fetchOidFromConfig(cfg SnmpPollingConfig, done chan SnmpFetchResult) {
	var result = SnmpFetchResult{Config: cfg}

	//time.Sleep(time.Duration(idx * 100000000))
	var snmpver gosnmp.SnmpVersion
	var msgflags gosnmp.SnmpV3MsgFlags
	var securityParams gosnmp.UsmSecurityParameters
	if cfg.SnmpVersion == "SNMP2c" {
		snmpver = gosnmp.Version2c
	} else if cfg.SnmpVersion == "SNMP1" {
		snmpver = gosnmp.Version1
	} else if cfg.SnmpVersion == "SNMP3" {
		snmpver = gosnmp.Version3

		if cfg.SnmpV3SecurityLevel == "authPriv" {
			msgflags = gosnmp.AuthPriv
		} else if cfg.SnmpV3SecurityLevel == "authNoPriv" {
			msgflags = gosnmp.AuthNoPriv
		} else {
			msgflags = gosnmp.NoAuthNoPriv
		}
		msgflags |= gosnmp.Reportable

		var authProtocol gosnmp.SnmpV3AuthProtocol
		if cfg.SnmpV3AuthenticationProtocol == "SHA" {
			authProtocol = gosnmp.SHA
		} else {
			authProtocol = gosnmp.MD5
		}

		var privProtocol gosnmp.SnmpV3PrivProtocol
		if cfg.SnmpV3PrivacyProtocol == "AES" {
			privProtocol = gosnmp.AES
		} else {
			privProtocol = gosnmp.DES
		}

		securityParams = gosnmp.UsmSecurityParameters{UserName: cfg.SnmpV3SecurityName,
			AuthenticationProtocol:   authProtocol,
			AuthenticationPassphrase: cfg.SnmpV3AuthenticationPassphrase,
			PrivacyProtocol:          privProtocol,
			PrivacyPassphrase:        cfg.SnmpV3PrivacyPassphrase,
		}
	}
	conn := &gosnmp.GoSNMP{
		Target:             cfg.IpAddress,
		Port:               161,
		Community:          cfg.SnmpCommunityName,
		Version:            snmpver,
		MsgFlags:           msgflags,
		SecurityModel:      gosnmp.UserSecurityModel,
		SecurityParameters: &securityParams,
		Timeout:            time.Duration(cfg.SnmpTimeout*cfg.SnmpRetries) * time.Second,
		Retries:            cfg.SnmpRetries,
		MaxRepetitions:     5,
	}

	result.Err = conn.Connect()
	if result.Err != nil {
		done <- result
		return
	}
	defer conn.Conn.Close()

	var data []gosnmp.SnmpPDU
	if cfg.PollType == "Walk" || cfg.PollType == "Table" {
		var res []gosnmp.SnmpPDU
		if conn.Version == gosnmp.Version1 {
			res, result.Err = conn.WalkAll(cfg.Oid)
		} else {
			res, result.Err = conn.BulkWalkAll(cfg.Oid)
		}
		if result.Err != nil {
			done <- result
			return
		}
		data = res
	} else if cfg.PollType == "Get" {
		var resp *gosnmp.SnmpPacket
		resp, result.Err = conn.Get([]string{cfg.Oid})
		if result.Err != nil {
			done <- result
			return
		}
		data = resp.Variables
	}
	result = updatePollTimes(result)
	result.Data = data
	done <- result
}

// update poll time fields in the snmpPollingConfig structure
func updatePollTimes(result SnmpFetchResult) (res SnmpFetchResult) {
	res = result
	// this time math is used to generate a poll time between the start of the next timeslot and 2 minutes before the next timeslot ends.
	current := time.Now()
	year, month, day := current.Date()
	today := time.Date(year, month, day, 0, 0, 0, 0, time.Local)
	freq := time.Duration(res.Config.PollFreq) * time.Second

	current_daily_timeslot := current.Sub(today) / freq
	next_timeslot_start := today.Add((current_daily_timeslot + 1) * freq)

	next_poll_start := next_timeslot_start.Add((time.Duration(rand.Intn(int(float64(res.Config.PollFreq)*0.8))) * time.Second) + (time.Duration(rand.Intn(1000)) * time.Millisecond))

	res.Config.LastPollTime = Now()
	res.Config.NextPollTime = next_poll_start.UnixNano() / int64(time.Millisecond)
	return
}

// update poll time fields in nmsConfigurationRemote.snmpPollingConfig table
func updateDbPollTimes(c SnmpPollingConfig, dbmap *gorp.DbMap) (err error) {
	var q = "" +
		"UPDATE `nmsConfigurationRemote`.`snmpPollingConfig`\n" +
		"SET `lastPollTime` = ?, `nextPollTime` = ?\n" +
		"WHERE resourceName = ? AND oid = ?"
	_, err = dbmap.Exec(q, c.LastPollTime, c.NextPollTime, c.ResourceName, c.Oid)
	return err
}
