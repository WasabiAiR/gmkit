# icinga

This is a tool kit to help interact with icinga.

You should be able to do the following today.
- Check to see if a hosts exists
- Enable and Disable Active checks
- Enable and Disable Notifications
- Set or Reset Downtime

## Flags and Environment
Currently here is the list of valid flags that can be used
```
  -icinga string
    	icinga address
  -icinga-password string
    	icinga password (optional)
  -icinga-tls-ca-cert string
    	The CA cert to use to validate the Icinga server's certificate
  -icinga-tls-client-cert string
    	The TLS client certificate to use when connecting to Icinga
  -icinga-tls-client-key string
    	The TLS client key to use when connecting to Icinga
  -icinga-tls-insecure
    	Whether or not to perform certificate hostname validation when connecting to the Nomad server
  -icinga-username string
    	icinga username (optional)
```

Currently the ENV this module will look for
```
gm_icinga
gm_icinga_username
gm_icinga_password
gm_icinga_tls_client_cert
gm_icinga_tls_client_key
gm_icinga_tls_ca_cert
gm_icinga_tls_insecure
```

## Connection
It is recommended to use the TLS certs that icinga installation addes to `/etc/icinga2/pki` and not use the username/password.

Be default nodes do not have access to the icinga api.  To give permission you need to add the following to the puppet config.
```
  @@firewall { "100 allow icinga2 api ${::fqdn}":
      state  => 'NEW',
      dport  => [5665],
      proto  => 'tcp',
      action => 'accept',
      source => $::facts[$ip_fact],
      tag    => ["api-${icinga::master_name}"]
    }

  # Configure an icinga api user for this machine
  @@icinga2::object::apiuser { "${::fqdn}-api" :
    client_cn   => $::fqdn,
    target      => "/etc/icinga2/zones.d/${icinga::master_name}/${::fqdn}.conf",
    permissions => [ '*' ],
    tag         => ["api-${icinga::master_name}"]
  }
```

## Setup Client
```
// Parse the flags
icingaCfg := icinga.NewConfig(flag.CommandLine)
flag.Parse()

// Create a icinga client
ic, err := icingaCfg.Client()
if err != nil {
  l.Fatal("error", err)
}
```

## Squelch a Host and all services
```
// Turn off notifications
if err := ic.SetAllNotifications(hostname, false); err != nil {
  return err
}

// Turn off active checks
if err := ic.SetAllActiveChecks(hostname, false); err != nil {
  return err
}

// Turn on downtimes
if err := ic.SetAllDowntime(hostname, "icinga-squelch", "Icinga squelch", time.Now(), time.Now().Add(time.Hour*24)); err != nil {
  return err
}
```
