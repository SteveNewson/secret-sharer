# Summary
Some tooling for sharing secrets over an insecure channel.

# Examples
Split a secret into 5 shares with a minimum threshold of 2:
```shell
> echo "thepassword" | secret-sharer split -n 5 -t 2
4a387c3781fab11b7648558317
a49bef88a955720ae51e1a2c0b
ca80312a83b4f8d13b3744cdb2
647f0916aa3c92a503d38e4571
1f66e063221970f0eac6e6601d
```

Receive two encrypted secrets (from two different senders), decrypt them and combine them:
```shell
> secret-sharer receive -s 2 | secret-sharer combine
```

The receive command will generate and print a transport key. That key needs to be used to wrap each sender's share:
```shell
> echo -n "4a387c3781fab11b7648558317" | secret-sharer wrap --transport-key "CLjQoM4K...."
.....aDgERguKn6mRlG2....v1p9c1M=
> echo -n "647f0916aa3c92a503d38e4571" | secret-sharer wrap --transport-key "CLjQoM4K...."
.....DgEvQbPyS6Us97R....Lbap7Qx=
```

The wrapped secret can then be transferred to the waiting receiver over an insecure channel, where they can be pasted into
the waiting request. 

Typically, the output of the combine command will be piped directly into another command, thereby avoiding the plaintext 
key from being serialised in any persistent form (or indeed the operator seeing the secret).