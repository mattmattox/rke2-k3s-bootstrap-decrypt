# RKE2/k3s Bootstrap decryption tool

This tool is used to decrypt the bootstrap key that stored in etcd by RKE2/k3s.

## Build

```bash
go build -o rke2-k3s-bootstrap-decryptor
```

## Usage

```bash
./rke2-k3s-bootstrap-decryptor -passphrase "1234567890"
```

To use this tool, you need to have the encryption key that used to encrypt the bootstrap key. This key is stored in the file `/var/lib/rancher/rke2/server/node-token` on the master nodes.

```bash
cat /var/lib/rancher/rke2/server/node-token
```

**NOTE**:The token is stored in the format `K10<CA-HASH>::<USERNAME>:<PASSWORD>`. We only need the password at the end.

We also need to capture the bootstrap key from etcd and store the output in a file called `bootstrap`

```bash
ETCD_POD=$(kubectl -n kube-system get pods -l component=etcd -o name | awk -F '/' '{print $2}' | head -n1)
kubectl -n kube-system exec -it ${ETCD_POD} -- sh
ETCDCTL_API=3 etcdctl --cert /var/lib/rancher/rke2/server/tls/etcd/server-client.crt --key /var/lib/rancher/rke2/server/tls/etcd/server-client.key --endpoints https://127.0.0.1:2379 --cacert /var/lib/rancher/rke2/server/tls/etcd/server-ca.crt get --prefix / --keys-only | grep "bootstrap/"
```

Example output:

```bash
/bootstrap/abcdef123456
```

Now that we have the bootstrap key name, we can run the following command to see the data stored in etcd:

```bash
ETCD_POD=$(kubectl -n kube-system get pods -l component=etcd -o name | awk -F '/' '{print $2}' | head -n1)
kubectl -n kube-system exec -it ${ETCD_POD} -- sh
ETCDCTL_API=3 etcdctl --cert /var/lib/rancher/rke2/server/tls/etcd/server-client.crt --key /var/lib/rancher/rke2/server/tls/etcd/server-client.key --endpoints https://127.0.0.1:2379 --cacert /var/lib/rancher/rke2/server/tls/etcd/server-ca.crt get /bootstrap/abcdef123456
```

Example output:

```bash
17b2c14a4de08362:mUrB4WbL/r1HSzuBP7i5SQfz52tC2tF6zdYu5e82qWDOjQBTxg0goakwlEe3h1eyT0SdXhfMKFjBSoIWo4f6fabQxKdPxlT407n1DPvVUE5BmWMEMYm+JswGlMfXvLxc6GJbzeOsIRO5jisf37MZryKXEPQBA7QGb8042PL4HDfMHHPN/tWc3+SUbZ0+Jbo8StUceZi1g1qwCKaHWf3Uk6Dc2uDqTadoeLD8ZGB98U6u9ROIV8i8DrTIzr3xZNNxM3sSDkhnJnfw7UT7/vWoH4/AjTqvjhBNblHuK/
...
```

**NOTE**: The output is the bootstrap key is an encrypted json string. We only need the value of the `key` field.

Now we can run the tool to decrypt the bootstrap key:

```bash
./rke2-k3s-bootstrap-decryptor -passphrase "1234567890"
```

Example YAML output:

```yaml
ClientCA:
  Content: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0t...
  Timestamp: "2023-10-10T21:03:26.236054153-05:00"
ClientCAKey:
  Content: LS0tLS1CRUdJTiBFQyBQUklWQVRFIEtFWS0r...
  Timestamp: "2023-10-10T21:03:26.236054153-05:00"
ETCDPeerCA:
  Content: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0r...
  Timestamp: "2023-10-10T21:03:26.264054702-05:00"
ETCDPeerCAKey:
  Content: LS0tLS1CRUdJTiBFQyBQUklWQVRFIEtFWS0r...
  Timestamp: "2023-10-10T21:03:26.264054702-05:00"
ETCDServerCA:
  Content: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0t...
  Timestamp: "2023-10-10T21:03:26.260054624-05:00"
ETCDServerCAKey:
  Content: LS0tLS1CRUdJTiBFQyBQUklWQVRFIEtFWS0r...
  Timestamp: "2023-10-10T21:03:26.260054624-05:00"
EncryptionConfig:
  Content: eyJraW5kIjoiRW5jcnlwdGlvbkNvbmZpZ3Vr...
  Timestamp: "2023-10-10T21:03:26.596061215-05:00"
EncryptionHash:
  Content: c3RhcnQtMDRlZTY0OTliMDIyMjI5ZTU0YTcz...
  Timestamp: "2023-10-10T21:03:26.596061215-05:00"
IPSECKey:
  Content: NmFmMjJlMjNlMWZmYjc5ZDI4YjY4NDhlZGNj...
  Timestamp: "2023-10-10T21:03:26.592061136-05:00"
PasswdFile:
  Content: Njc3YjExMTg0YTRmNWFhODBkNWUzMDg2ZDI1...
  Timestamp: "2023-10-10T21:03:26.592061136-05:00"
RequestHeaderCA:
  Content: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0t...
  Timestamp: "2023-10-10T21:03:26.260054624-05:00"
RequestHeaderCAKey:
  Content: LS0tLS1CRUdJTiBFQyBQUklWQVRFIEtFWS0t...
  Timestamp: "2023-10-10T21:03:26.256054546-05:00"
ServerCA:
  Content: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0t...
  Timestamp: "2023-10-10T21:03:26.256054546-05:00"
ServerCAKey:
  Content: LS0tLS1CRUdJTiBFQyBQUklWQVRFIEtFWS0t...
  Timestamp: "2023-10-10T21:03:26.252054467-05:00"
ServiceKey:
  Content: LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVkt...
  Timestamp: "2023-10-10T21:03:26.592061136-05:00"
```

**NOTE**: Each of Content field is a base64 encoded string. You can decode it using the following command:

```bash
echo -n "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0t" | base64 -d; echo
```

As you can see, the output contains all the certificates and keys that are used by RKE2/k3s. Which is why it is important to keep the token safe and secure as with it you can own the cluster.

## License

This project is licensed under the Apache-2.0 License.