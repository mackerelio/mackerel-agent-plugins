# Changelog

## 0.88.1 (2025-03-31)

* replace to newer runner-images #1264 (yseto)
* Bump golang.org/x/net from 0.33.0 to 0.36.0 in /mackerel-plugin-gcp-compute-engine #1263 (dependabot[bot])


## 0.88.0 (2025-03-04)

* support memory peak value which is introduced in PHP 8.4 on mackerel-plugin-php-fpm #1254 (kmuto)
* Bump mackerelio/workflows from 1.2.0 to 1.3.0 #1251 (dependabot[bot])


## 0.87.0 (2025-01-27)

* Fix CI build #1247 (ne-sachirou)
* Bump golang.org/x/net from 0.23.0 to 0.33.0 in /mackerel-plugin-gcp-compute-engine #1246 (dependabot[bot])
* Bump golang.org/x/crypto from 0.17.0 to 0.31.0 in /mackerel-plugin-mongodb #1244 (dependabot[bot])
* use mackerelio/workflows@v1.2.0 #1243 (yseto)
* Bump golang.org/x/text from 0.17.0 to 0.21.0 #1242 (dependabot[bot])
* Bump golang.org/x/sync from 0.8.0 to 0.10.0 #1241 (dependabot[bot])
* Bump github.com/stretchr/testify from 1.8.1 to 1.10.0 in /mackerel-plugin-nvidia-smi #1240 (dependabot[bot])
* Bump github.com/stretchr/testify from 1.9.0 to 1.10.0 #1239 (dependabot[bot])
* [jvm] add prefix option #1238 (masarasi)
* Bump github.com/redis/go-redis/v9 from 9.5.1 to 9.7.0 #1233 (dependabot[bot])
* Bump github.com/urfave/cli from 1.22.15 to 1.22.16 #1232 (dependabot[bot])
* Bump github.com/mackerelio/go-mackerel-plugin-helper from 0.1.2 to 0.1.3 in /mackerel-plugin-nvidia-smi #1230 (dependabot[bot])
* Bump github.com/mackerelio/mackerel-plugin-mongodb from 1.1.0 to 1.1.1 in /mackerel-plugin-mongodb #1227 (dependabot[bot])
* Bump github.com/mackerelio/mackerel-plugin-json from 1.2.4 to 1.2.5 in /mackerel-plugin-json #1226 (dependabot[bot])
* Bump github.com/opencontainers/runc from 1.1.12 to 1.1.14 #1221 (dependabot[bot])
* Bump github.com/gosnmp/gosnmp from 1.35.0 to 1.38.0 #1215 (dependabot[bot])
* Bump github.com/aws/aws-sdk-go from 1.51.26 to 1.55.5 #1212 (dependabot[bot])


## 0.86.0 (2024-10-08)

* Bump github.com/mackerelio/mackerel-plugin-mysql from 1.2.1 to 1.3.0 in /mackerel-plugin-mysql #1228 (dependabot[bot])
* Watch every module go.mod #1225 (yohfee)
* Bump github.com/mackerelio/mackerel-plugin-mysql from 1.2.2 to 1.3.0 #1224 (dependabot[bot])
* Fix lint #1220 (yohfee)
* [plugin-aws-ec2-ebs] add actions to README need to be allowed in the iam policy #1219 (miztch)


## 0.85.0 (2024-08-08)

* Bump golang.org/x/text from 0.15.0 to 0.17.0 #1216 (dependabot[bot])
* [mackerel-plugin-snmp] support 32bit counter overflow. #1214 (yseto)
* Bump github.com/mackerelio/go-mackerel-plugin from 0.1.4 to 0.1.5 #1198 (dependabot[bot])
* Bump github.com/mackerelio/go-mackerel-plugin-helper from 0.1.2 to 0.1.3 #1196 (dependabot[bot])
* Bump github.com/mackerelio/go-osstat from 0.2.4 to 0.2.5 #1195 (dependabot[bot])
* Bump github.com/jmoiron/sqlx from 1.3.5 to 1.4.0 #1178 (dependabot[bot])


## 0.84.0 (2024-07-01)

* [mackerel-plugin-php-fpm] add slow_requests_delta metrics because slow_requests is a counter #1201 (Arthur1)


## 0.83.0 (2024-06-12)

* [plugin-mailq] mailq/postfix: Fix mail queue count regexp capture when it's `1` #1186 (yongjiajun)
* Bump github.com/urfave/cli from 1.22.14 to 1.22.15 #1181 (dependabot[bot])
* Bump golang.org/x/sync from 0.6.0 to 0.7.0 #1180 (dependabot[bot])
* Bump github.com/stretchr/testify from 1.8.4 to 1.9.0 #1177 (dependabot[bot])


## 0.82.1 (2024-04-23)

* Bump github.com/aws/aws-sdk-go from 1.45.11 to 1.51.26 #1175 (dependabot[bot])
* Bump golang.org/x/net from 0.7.0 to 0.23.0 in /mackerel-plugin-gcp-compute-engine #1172 (dependabot[bot])
* Bump github.com/docker/docker from 23.0.1+incompatible to 24.0.9+incompatible #1163 (dependabot[bot])
* Bump google.golang.org/protobuf from 1.28.1 to 1.33.0 in /mackerel-plugin-gcp-compute-engine #1160 (dependabot[bot])
* Bump google.golang.org/protobuf from 1.27.1 to 1.33.0 in /mackerel-plugin-murmur #1159 (dependabot[bot])
* Bump github.com/redis/go-redis/v9 from 9.1.0 to 9.5.1 #1148 (dependabot[bot])
* Bump github.com/opencontainers/runc from 1.1.2 to 1.1.12 #1147 (dependabot[bot])
* Bump github.com/containerd/containerd from 1.6.18 to 1.6.26 #1143 (dependabot[bot])
* Bump golang.org/x/crypto from 0.0.0-20220622213112-05595931fe9d to 0.17.0 in /mackerel-plugin-mongodb #1142 (dependabot[bot])
* Bump google.golang.org/grpc from 1.51.0 to 1.56.3 in /mackerel-plugin-gcp-compute-engine #1125 (dependabot[bot])
* Bump github.com/go-redis/redismock/v9 from 9.0.3 to 9.2.0 #1122 (dependabot[bot])
* Bump github.com/jarcoal/httpmock from 1.3.0 to 1.3.1 #1111 (dependabot[bot])


## 0.82.0 (2024-04-05)

* add STRING number support #1166 (kmuto)
* Allow mackerel-plugin-redis to specify username #1165 (mkadokawa-idcf)


## 0.81.0 (2024-03-05)

* [mackerel-plugin-windows-server-sessions] Remove the dependency for WMIC command #1155 (Arthur1)


## 0.80.0 (2024-02-27)

* Update mysql, mongodb plugins. #1152 (yseto)
* update go version -> 1.22 #1145 (lufia)
* added TLS support on mackerel-plugin-redis #1144 (yseto)
* Bump golang.org/x/crypto from 0.0.0-20220622213112-05595931fe9d to 0.17.0 #1141 (dependabot[bot])
* Bump actions/upload-artifact from 3 to 4 #1140 (dependabot[bot])
* Bump actions/download-artifact from 3 to 4 #1139 (dependabot[bot])
* Bump actions/setup-go from 4 to 5 #1136 (dependabot[bot])


## 0.79.0 (2023-09-22)

* Bump github.com/aws/aws-sdk-go from 1.44.239 to 1.45.11 #1109 (dependabot[bot])
* Bump actions/checkout from 3 to 4 #1108 (dependabot[bot])
* use go-redis #1105 (yseto)
* Bump golang.org/x/text from 0.9.0 to 0.13.0 #1103 (dependabot[bot])
* Bump golang.org/x/sync from 0.1.0 to 0.3.0 #1082 (dependabot[bot])
* Bump github.com/urfave/cli from 1.22.12 to 1.22.14 #1080 (dependabot[bot])
* Bump github.com/montanaflynn/stats from 0.7.0 to 0.7.1 #1068 (dependabot[bot])
* Bump github.com/lib/pq from 1.10.7 to 1.10.9 #1063 (dependabot[bot])


## 0.78.4 (2023-09-04)

* Fixed Docker CPU Percentage's unusual spike when restarting containers #1100 (Arthur1)
* Remove old rpm packaging #1095 (yseto)


## 0.78.3 (2023-07-13)



## 0.78.2 (2023-06-14)

* Bump github.com/mackerelio/mackerel-plugin-mysql from 1.2.0 to 1.2.1 in /mackerel-plugin-mysql #1077 (dependabot[bot])
* Bump github.com/mackerelio/mackerel-plugin-mysql from 1.2.0 to 1.2.1 #1076 (dependabot[bot])
* [fluentd] add README.md with extended_metrics option #1072 (sakamossan)


## 0.78.1 (2023-04-12)

* Bump golang.org/x/text from 0.8.0 to 0.9.0 #1057 (dependabot[bot])
* Bump github.com/aws/aws-sdk-go from 1.44.229 to 1.44.239 #1056 (dependabot[bot])
* Fix timeout/error labels of mackerel-plugin-twemproxy #1055 (kmuto)
* [ci] refactor .github/workflows/test.yml #1050 (lufia)
* Bump github.com/aws/aws-sdk-go from 1.44.219 to 1.44.229 #1049 (dependabot[bot])
* Bump github.com/mackerelio/go-osstat from 0.2.3 to 0.2.4 #1046 (dependabot[bot])


## 0.78.0 (2023-03-13)

* Bump github.com/fsouza/go-dockerclient from 1.9.5 to 1.9.6 #1044 (dependabot[bot])
* Bump github.com/aws/aws-sdk-go from 1.44.199 to 1.44.219 #1043 (dependabot[bot])
* Bump golang.org/x/text from 0.7.0 to 0.8.0 #1042 (dependabot[bot])
* Bump github.com/mackerelio/mackerel-plugin-mongodb from 1.0.0 to 1.1.0 in /mackerel-plugin-mongodb #1037 (dependabot[bot])
* Bump github.com/mackerelio/mackerel-plugin-mongodb from 1.0.0 to 1.1.0 #1033 (dependabot[bot])
* Fix support APCu and PHP7 or later environment. #1032 (uzulla)


## 0.77.0 (2023-02-27)

* Bump github.com/mackerelio/mackerel-plugin-mysql from 1.1.0 to 1.2.0 #1039 (dependabot[bot])
* Bump github.com/mackerelio/mackerel-plugin-mysql from 1.1.0 to 1.2.0 in /mackerel-plugin-mysql #1038 (dependabot[bot])
* Bump github.com/stretchr/testify from 1.8.1 to 1.8.2 #1036 (dependabot[bot])
* Bump github.com/fsouza/go-dockerclient from 1.9.4 to 1.9.5 #1035 (dependabot[bot])
* fix syntax for aws rds #1031 (heleeen)
* Bump golang.org/x/net from 0.3.0 to 0.7.0 in /mackerel-plugin-gcp-compute-engine #1029 (dependabot[bot])
* Bump github.com/containerd/containerd from 1.6.14 to 1.6.18 #1028 (dependabot[bot])
* added multiple-os tests #1026 (yseto)
* added import. #1025 (yseto)
* added external repository mongodb #1023 (yseto)


## 0.76.1 (2023-02-15)

* Bump golang.org/x/text from 0.6.0 to 0.7.0 #1022 (dependabot[bot])
* Bump github.com/fsouza/go-dockerclient from 1.9.3 to 1.9.4 #1021 (dependabot[bot])
* Bump github.com/aws/aws-sdk-go from 1.44.190 to 1.44.199 #1020 (dependabot[bot])
* Remove `circle.yml` #1019 (wafuwafu13)
* plugin-mongodb is external repository #1018 (yseto)
* Bump github.com/jarcoal/httpmock from 1.2.0 to 1.3.0 #1016 (dependabot[bot])


## 0.76.0 (2023-02-01)

* Bump github.com/aws/aws-sdk-go from 1.44.184 to 1.44.190 #1014 (dependabot[bot])
* Bump actions/setup-go from 2 to 3 #1013 (dependabot[bot])
* Bump actions/cache from 2 to 3 #1012 (dependabot[bot])
* Bump actions/upload-artifact from 2 to 3 #1011 (dependabot[bot])
* Bump actions/download-artifact from 2 to 3 #1010 (dependabot[bot])
* Bump actions/checkout from 2 to 3 #1009 (dependabot[bot])
* Bump github.com/mackerelio/mackerel-plugin-mysql from 1.0.0 to 1.1.0 in /mackerel-plugin-mysql #1008 (dependabot[bot])
* Bump github.com/urfave/cli from 1.22.10 to 1.22.12 #1006 (dependabot[bot])
* Bump github.com/mackerelio/mackerel-plugin-mysql from 1.0.0 to 1.1.0 #1005 (dependabot[bot])
* Enables Dependabot version updates for GitHub Actions #1004 (Arthur1)
* Remove debian package v1 process. #1003 (yseto)
* fix staticcheck #1001 (yseto)
* Bump github.com/fsouza/go-dockerclient from 1.9.0 to 1.9.3 #999 (dependabot[bot])
* Bump github.com/aws/aws-sdk-go from 1.44.116 to 1.44.184 #998 (dependabot[bot])
* fix gosimple #997 (yseto)
* fix errcheck, ineffassign. #996 (yseto)
* ci: enable `gofmt` #994 (wafuwafu13)
* Bump golang.org/x/text from 0.5.0 to 0.6.0 #991 (dependabot[bot])
* Bump github.com/montanaflynn/stats from 0.6.6 to 0.7.0 #990 (dependabot[bot])


## 0.75.0 (2023-01-18)

* fix build on current working directory. #985 (yseto)
* added compile option. #983 (yseto)
* plugin-mysql is external repository #982 (yseto)
* packaging: make compression format xz #981 (lufia)
* accesslog: use `reqtime_microsec` if exists #976 (wafuwafu13)
* split go.mod for plugins that have previously split repositories #972 (yseto)


## 0.74.0 (2022-12-20)

* use ubuntu-20.04 #978 (yseto)
* fix packaging process on ci #971 (yseto)
* refine file rewrite process #970 (yseto)
* remove xentop on backward compatibility symlink #969 (yseto)
* fix packaging #968 (yseto)
* [plugin-elasticsearch] Fix the test for elasticsearch #966 (lufia)
* added external plugin support #964 (yseto)
* sort plugins on packaging files. #963 (yseto)
* Purge less-used plugins from mackerel-agent-plugins package #962 (lufia)
* Purge mackerel-plugin-nvidia-smi #961 (lufia)
* Update dependencies #960 (lufia)
* added test for elasticsearch #952 (yseto)


## 0.73.0 (2022-11-09)

* Fix Elasticsearch plugin. #950 (fujiwara)
* Bump github.com/stretchr/testify from 1.8.0 to 1.8.1 #946 (dependabot[bot])
* [mongodb] add metric-key-prefix option #945 (tukaelu)


## 0.72.2 (2022-10-20)

* Bump github.com/aws/aws-sdk-go from 1.44.56 to 1.44.116 #942 (dependabot[bot])
* Bump github.com/fsouza/go-dockerclient from 1.8.3 to 1.9.0 #941 (dependabot[bot])
* added timeout on lint #935 (yseto)
* Bump github.com/urfave/cli from 1.22.9 to 1.22.10 #933 (dependabot[bot])
* Bump github.com/lib/pq from 1.10.5 to 1.10.7 #931 (dependabot[bot])
* [uptime] Add tests #929 (wafuwafu13)
* go.mod from 1.16 to 1.18 #928 (yseto)
* Bump github.com/mackerelio/go-mackerel-plugin from 0.1.3 to 0.1.4 #927 (dependabot[bot])
* Bump github.com/fsouza/go-dockerclient from 1.8.1 to 1.8.3 #925 (dependabot[bot])
* Bump github.com/mackerelio/go-osstat from 0.2.2 to 0.2.3 #924 (dependabot[bot])
* Improve test #923 (yseto)
* Bump github.com/go-ldap/ldap/v3 from 3.4.3 to 3.4.4 #917 (dependabot[bot])
* [aws-ec2-ebs] fix calcurate procedure of Nitro instance #916 (yseto)
* [plugin-aws-ec2-ebs] fix misuse of period #915 (lufia)
* Bump github.com/mackerelio/go-mackerel-plugin-helper from 0.1.1 to 0.1.2 #914 (dependabot[bot])
* Bump github.com/gosnmp/gosnmp from 1.34.0 to 1.35.0 #889 (dependabot[bot])
* Bump github.com/jarcoal/httpmock from 1.1.0 to 1.2.0 #883 (dependabot[bot])


## 0.72.1 (2022-07-20)

* Bump github.com/aws/aws-sdk-go from 1.44.37 to 1.44.56 #910 (dependabot[bot])
* Bump github.com/gomodule/redigo from 1.8.8 to 1.8.9 #908 (dependabot[bot])
* Bump github.com/stretchr/testify from 1.7.1 to 1.8.0 #906 (dependabot[bot])
* Bump github.com/hashicorp/go-version from 1.4.0 to 1.6.0 #905 (dependabot[bot])
* Bump github.com/mackerelio/go-mackerel-plugin from 0.1.2 to 0.1.3 #902 (dependabot[bot])
* Bump github.com/urfave/cli from 1.22.5 to 1.22.9 #881 (dependabot[bot])
* Bump github.com/jmoiron/sqlx from 1.3.4 to 1.3.5 #872 (dependabot[bot])


## 0.72.0 (2022-06-22)

* Bump github.com/aws/aws-sdk-go from 1.43.36 to 1.44.37 #900 (dependabot[bot])
* Bump github.com/fsouza/go-dockerclient from 1.8.0 to 1.8.1 #899 (dependabot[bot])
* [plugin-docker] update README and a description of -command option #896 (lufia)
* [plugin-docker] fix CPU/Memory metrics on Docker hosts uses cgroup2 #895 (xruins)
* [plugin-docker] (Breaking) drop 'File' method support #894 (lufia)
* Bump github.com/fsouza/go-dockerclient from 1.7.10 to 1.8.0 #892 (dependabot[bot])


## 0.71.0 (2022-04-21)

* [plugin-mysql] fix panic when parsing aio stats #874 (lufia)
* Fix: Input 'job-number' has been deprecated with message: use flag-name instead #871 (ne-sachirou)


## 0.70.6 (2022-04-14)

* Bump github.com/go-ldap/ldap/v3 from 3.4.2 to 3.4.3 #869 (dependabot[bot])
* Bump github.com/aws/aws-sdk-go from 1.43.26 to 1.43.36 #868 (dependabot[bot])
* Bump github.com/lib/pq from 1.10.4 to 1.10.5 #867 (dependabot[bot])
* [linux] users メトリックに 0 が投稿されない問題を修正 #866 (masarasi)


## 0.70.5 (2022-03-30)

* Bump github.com/mackerelio/go-osstat from 0.2.1 to 0.2.2 #863 (dependabot[bot])
* Bump github.com/aws/aws-sdk-go from 1.43.11 to 1.43.26 #862 (dependabot[bot])
* Bump github.com/stretchr/testify from 1.7.0 to 1.7.1 #859 (dependabot[bot])
* Bump github.com/fsouza/go-dockerclient from 1.7.9 to 1.7.10 #856 (dependabot[bot])


## 0.70.4 (2022-03-15)

* Bump github.com/aws/aws-sdk-go from 1.42.52 to 1.43.11 #854 (dependabot[bot])
* Bump github.com/fsouza/go-dockerclient from 1.7.8 to 1.7.9 #853 (dependabot[bot])
* Bump github.com/go-ldap/ldap/v3 from 3.4.1 to 3.4.2 #851 (dependabot[bot])


## 0.70.3 (2022-02-16)

* Bump github.com/aws/aws-sdk-go from 1.40.59 to 1.42.52 #848 (dependabot[bot])
* Bump github.com/fsouza/go-dockerclient from 1.7.4 to 1.7.8 #847 (dependabot[bot])
* Bump github.com/gomodule/redigo from 1.8.6 to 1.8.8 #839 (dependabot[bot])
* Bump github.com/hashicorp/go-version from 1.3.0 to 1.4.0 #838 (dependabot[bot])


## 0.70.2 (2022-01-12)

* Bump github.com/jarcoal/httpmock from 1.0.8 to 1.1.0 #836 (dependabot[bot])
* Bump github.com/gomodule/redigo from 1.8.5 to 1.8.6 #832 (dependabot[bot])
* Bump github.com/gosnmp/gosnmp from 1.32.0 to 1.34.0 #830 (dependabot[bot])
* Bump github.com/lib/pq from 1.10.3 to 1.10.4 #826 (dependabot[bot])


## 0.70.1 (2021-12-01)

* upgrade to Go 1.17 and others #828 (lufia)
* Bump github.com/mackerelio/go-osstat from 0.2.0 to 0.2.1 #818 (dependabot[bot])


## 0.70.0 (2021-11-18)

* [plugin-sidekiq] add queue latency metric #823 (ch1aki)


## 0.69.1 (2021-10-14)

* Bump github.com/aws/aws-sdk-go from 1.39.4 to 1.40.59 #816 (dependabot[bot])
* Bump github.com/fsouza/go-dockerclient from 1.7.3 to 1.7.4 #814 (dependabot[bot])
* Bump github.com/go-ldap/ldap/v3 from 3.3.0 to 3.4.1 #813 (dependabot[bot])
* Bump github.com/lib/pq from 1.10.2 to 1.10.3 #812 (dependabot[bot])
* Bump github.com/mackerelio/golib from 1.2.0 to 1.2.1 #811 (dependabot[bot])


## 0.69.0 (2021-09-29)

* [plugin-redis] migrate redis client library to redigo #809 (pyto86pri)
* replace library gosnmp in mackerel-plugin-snmp #808 (yseto)
* [mysql] adapt for mysql5.7's InnoDB SEMAPHORES metrics #807 (do-su-0805)


## 0.68.0 (2021-08-24)

* [plugin-postgres] suppress fetchXlogLocation when wal_level != 'logical' #804 (handlename)
* added "pending log flushes" #802 (yseto)
* fix pending_normal_aio_* #801 (yseto)


## 0.67.1 (2021-08-05)

* [accesslog] use Seek to skip the log #795 (lufia)
* [ci][fluentd] add test.sh #794 (lufia)


## 0.67.0 (2021-07-28)

* [fluentd] add workers option. #791 (sugy)
* Bump github.com/aws/aws-sdk-go from 1.38.70 to 1.39.4 #789 (dependabot[bot])


## 0.66.0 (2021-07-15)

* [fluentd] add metric-key-prefix option. #788 (sugy)


## 0.65.0 (2021-07-06)

* bump go-mackerel-plugin and go-mackerel-plugin-helper #783 (astj)
* Bump github.com/aws/aws-sdk-go from 1.38.40 to 1.38.70 #781 (dependabot[bot])
* Bump github.com/gomodule/redigo from 1.8.4 to 1.8.5 #773 (dependabot[bot])
* Bump github.com/fsouza/go-dockerclient from 1.7.2 to 1.7.3 #767 (dependabot[bot])
* Bump github.com/jmoiron/sqlx from 1.3.1 to 1.3.4 #761 (dependabot[bot])
* [ci][plugin-mysql] check metrics #778 (lufia)
* [plugin-aws-cloudfront]Replace label name of graph #766 (wafuwafu13)
* [plugin-redis] avoid to store +Inf #779 (lufia)


## 0.64.3 (2021-06-23)

* [ci][plugin-redis] check metrics #775 (lufia)
* [ci] run tests on the workflow #771 (lufia)


## 0.64.2 (2021-06-10)

* [plugin-mysql] Fix plugin-mysql to be able to collect metrics by ignoring non-numerical values #770 (lufia)


## 0.64.1 (2021-06-03)

* Bump github.com/montanaflynn/stats from 0.6.5 to 0.6.6 #758 (dependabot[bot])
* Bump github.com/lib/pq from 1.10.0 to 1.10.2 #762 (dependabot[bot])
* Bump github.com/aws/aws-sdk-go from 1.38.1 to 1.38.40 #760 (dependabot[bot])
* Bump github.com/mackerelio/go-osstat from 0.1.0 to 0.2.0 #759 (dependabot[bot])
* Bump github.com/go-ldap/ldap/v3 from 3.2.4 to 3.3.0 #750 (dependabot[bot])
* Use latest go-mackerel-plugin(-helper) #756 (astj)
* Bump github.com/hashicorp/go-version from 1.2.1 to 1.3.0 #748 (dependabot[bot])
* upgrade Go 1.14 -> 1.16 #753 (lufia)


## 0.64.0 (2021-03-24)

* Bump github.com/jmoiron/sqlx from 1.2.0 to 1.3.1 #722 (dependabot[bot])
* Bump github.com/aws/aws-sdk-go from 1.37.33 to 1.38.1 #744 (dependabot[bot])
* Bump github.com/aws/aws-sdk-go from 1.36.28 to 1.37.33 #743 (dependabot[bot])
* Bump github.com/gomodule/redigo from 1.8.3 to 1.8.4 #732 (dependabot[bot])
* Bump github.com/fsouza/go-dockerclient from 1.7.0 to 1.7.2 #739 (dependabot[bot])
* Bump github.com/lib/pq from 1.9.0 to 1.10.0 #742 (dependabot[bot])
* Bump github.com/montanaflynn/stats from 0.6.4 to 0.6.5 #734 (dependabot[bot])
* Bump github.com/mackerelio/golib from 1.1.0 to 1.2.0 #735 (dependabot[bot])
* Bump github.com/jarcoal/httpmock from 1.0.7 to 1.0.8 #721 (dependabot[bot])
* Support to extended metrics for fluentd 1.6. #738 (fujiwara)


## 0.63.5 (2021-02-19)

* fix incorrect endpoint inspect condition #730 (yseto)


## 0.63.4 (2021-02-19)

* migrate from goamz,go-ses to aws-sdk-go #726 (yseto)
* replace mackerel-github-release #724 (yseto)
* [plugin-redis] add test.sh #718 (lufia)


## 0.63.3 (2021-01-21)

* Revert "delete unused Makefile tasks" #716 (lufia)
* Bump github.com/aws/aws-sdk-go from 1.35.33 to 1.36.28 #712 (dependabot[bot])
* Bump github.com/fsouza/go-dockerclient from 1.6.6 to 1.7.0 #711 (dependabot[bot])
* Bump github.com/mackerelio/mackerel-plugin-json from 1.2.0 to 1.2.2 #706 (dependabot[bot])
* Bump github.com/montanaflynn/stats from 0.6.3 to 0.6.4 #713 (dependabot[bot])
* Bump github.com/jarcoal/httpmock from 1.0.6 to 1.0.7 #709 (dependabot[bot])
* Bump github.com/stretchr/testify from 1.6.1 to 1.7.0 #714 (dependabot[bot])
* Bump github.com/mackerelio/golib from 1.0.0 to 1.1.0 #715 (dependabot[bot])
* Fix labels of innodb_tables_in_use and innodb_locked_tables #708 (syou6162)
* Bump github.com/fsouza/go-dockerclient from 1.6.5 to 1.6.6 #696 (dependabot[bot])
* migrate garyburd/redigo -> gomodule/redigo #704 (lufia)
* migrate CI to GitHub Actions #701 (lufia)
* [plugin-apache2] add test.sh #700 (lufia)


## 0.63.2 (2020-12-09)

* Bump github.com/lib/pq from 1.8.0 to 1.9.0 #697 (dependabot[bot])
* Bump github.com/mackerelio/mackerel-plugin-json from 1.1.0 to 1.2.0 #692 (dependabot[bot])
* Enables to parse responses which include string and number from Plack. #694 (fujiwara)
* Bump github.com/aws/aws-sdk-go from 1.34.32 to 1.35.33 #691 (dependabot[bot])
* Bump github.com/urfave/cli from 1.22.4 to 1.22.5 #686 (dependabot-preview[bot])
* Bump github.com/go-ldap/ldap/v3 from 3.2.3 to 3.2.4 #682 (dependabot-preview[bot])
* Update Dependabot config file #689 (dependabot-preview[bot])


## 0.63.1 (2020-10-01)

* Bump github.com/aws/aws-sdk-go from 1.34.22 to 1.34.32 #677 (dependabot-preview[bot])
* fix build arch for make rpm-v2-arm #675 (astj)


## 0.63.0 (2020-09-15)

* add arm64 architecture packages, and fix Architecture field of deb #673 (lufia)
* Bump github.com/aws/aws-sdk-go from 1.33.17 to 1.34.22 #672 (dependabot-preview[bot])
* Bump github.com/jarcoal/httpmock from 1.0.5 to 1.0.6 #668 (dependabot-preview[bot])
* Bump github.com/go-redis/redis from 6.15.8+incompatible to 6.15.9+incompatible #664 (dependabot-preview[bot])
* Build with Go 1.14 #665 (lufia)
* Bump github.com/aws/aws-sdk-go from 1.33.12 to 1.33.17 #661 (dependabot-preview[bot])
* Bump github.com/garyburd/redigo from 1.6.0 to 1.6.2 #660 (dependabot-preview[bot])
* Bump github.com/lib/pq from 1.7.1 to 1.8.0 #662 (dependabot-preview[bot])
* Bump github.com/go-ldap/ldap/v3 from 3.1.10 to 3.2.3 #652 (dependabot-preview[bot])
* [plugin-mongodb] add test.sh #654 (lufia)


## 0.62.0 (2020-07-29)

* [plugin-linux] Do not skip diskstat metrics with Linux kernel 4.18+ #658 (astj)
* Bump github.com/aws/aws-sdk-go from 1.31.11 to 1.33.12 #656 (dependabot-preview[bot])
* Bump github.com/lib/pq from 1.7.0 to 1.7.1 #657 (dependabot-preview[bot])


## 0.61.0 (2020-07-20)

* [plugin-mysql] Fix to send Bytes_sent and Bytes_received correctly #649 (shibayu36)
* Bump github.com/lib/pq from 1.6.0 to 1.7.0 #641 (dependabot-preview[bot])
* Bump github.com/hashicorp/go-version from 1.2.0 to 1.2.1 #642 (dependabot-preview[bot])
* [plugin-accesslog] allow fields that cannot be parsed but unused by the plugin #645 (susisu)
* Bump github.com/go-redis/redis from 6.15.7+incompatible to 6.15.8+incompatible #631 (dependabot-preview[bot])
* Bump github.com/stretchr/testify from 1.6.0 to 1.6.1 #638 (dependabot-preview[bot])
* Update aws-sdk-go to 1.31.11 #636 (astj)
* Bump github.com/aws/aws-sdk-go from 1.30.27 to 1.31.7 #634 (dependabot-preview[bot])
* Bump github.com/lib/pq from 1.5.2 to 1.6.0 #633 (dependabot-preview[bot])
* Bump github.com/stretchr/testify from 1.5.1 to 1.6.0 #635 (dependabot-preview[bot])
* [plugin-postgres] add test.sh #629 (lufia)
* Bump github.com/aws/aws-sdk-go from 1.30.7 to 1.30.27 #627 (dependabot-preview[bot])


## 0.60.2 (2020-05-14)

* Bump github.com/fsouza/go-dockerclient from 1.6.3 to 1.6.5 #619 (dependabot-preview[bot])
* Bump github.com/lib/pq from 1.3.0 to 1.5.2 #624 (dependabot-preview[bot])
* Bump github.com/go-ldap/ldap/v3 from 3.1.8 to 3.1.10 #623 (dependabot-preview[bot])
* Bump github.com/Songmu/axslogparser from 1.2.0 to 1.3.0 #615 (dependabot-preview[bot])
* ignore diffs of go.mod and go.sum #626 (lufia)
* Bump go-mackerel-plugin{,-helper} #613 (astj)
* Bump github.com/aws/aws-sdk-go from 1.29.24 to 1.30.7 #612 (dependabot-preview[bot])
* Bump github.com/urfave/cli from 1.22.2 to 1.22.4 #609 (dependabot-preview[bot])
* Bump github.com/go-ldap/ldap/v3 from 3.1.7 to 3.1.8 #608 (dependabot-preview[bot])
* Add documents for testing #610 (lufia)


## 0.60.1 (2020-04-03)

* Bump github.com/jarcoal/httpmock from 1.0.4 to 1.0.5 #603 (dependabot-preview[bot])
* Bump github.com/aws/aws-sdk-go from 1.29.14 to 1.29.24 #601 (dependabot-preview[bot])
* Bump github.com/aws/aws-sdk-go from 1.28.13 to 1.29.14 #598 (dependabot-preview[bot])
* Bump github.com/montanaflynn/stats from 0.5.0 to 0.6.3 #596 (dependabot-preview[bot])
* Bump github.com/fsouza/go-dockerclient from 1.6.0 to 1.6.3 #595 (dependabot-preview[bot])
* Bump github.com/stretchr/testify from 1.4.0 to 1.5.1 #594 (dependabot-preview[bot])
* Bump github.com/go-ldap/ldap/v3 from 3.1.5 to 3.1.7 #591 (dependabot-preview[bot])
* Bump github.com/go-redis/redis from 6.15.6+incompatible to 6.15.7+incompatible #585 (dependabot-preview[bot])
* Bump github.com/aws/aws-sdk-go from 1.28.5 to 1.28.13 #588 (dependabot-preview[bot])


## 0.60.0 (2020-02-05)

* [varnish] remove unnecessary printf #586 (lufia)
* [varnish] Added metrics for Transient storage #584 (cohalz)
* [varnish] Added backend_reuse and backend_recycle metrics #583 (cohalz)
* rename: github.com/motemen/gobump -> github.com/x-motemen/gobump #582 (lufia)


## 0.59.1 (2020-01-22)

* Bump github.com/aws/aws-sdk-go from 1.27.0 to 1.28.5 #579 (dependabot-preview[bot])
* Bump github.com/aws/aws-sdk-go from 1.26.6 to 1.27.0 #577 (dependabot-preview[bot])
* add .dependabot/config.yml #575 (lufia)
* refactor Makefile and update dependencies #574 (lufia)


## 0.59.0 (2019-10-24)

* Build with Go 1.12.12
* Update dependencies #572 (astj)
* [solr] Fix several graph definitions #571 (supercaracal)
* [doc]add repository policy #570 (lufia)
* [jvm] Add the PerfDisableSharedMem JVM option issue to README #569 (supercaracal)


## 0.58.0 (2019-08-29)

* [solr] Add version 7.x and 8.x support #562 (supercaracal)
* [sidekiq] add option for redis-namespace #563 (y-imaida)
* add fakeroot to build dependencies #564 (susisu)


## 0.57.0 (2019-07-22)

* [jvm] fix jinfo command timed out error is not logged #560 (susisu)
* Build with Go 1.12 #558 (astj)
*  [plugin-jvm] added CGC and CGCT metrics and fixed parsing problem on them #557 (lufia)


## 0.56.0 (2019-06-11)

* support go modules #554 (astj)
* [plugin-jvm] prefer ${JAVA_HOME}/bin/j** if JAVA_HOME is set #553 (astj)


## 0.55.2 (2019-05-08)

* [mysql] add -debug option for troubleshooting #551 (lufia)


## 0.55.1 (2019-03-27)

* [accesslog] don't return any metrics on the first scan #549 (Songmu)
* [haproxy] fix example of haproxy.cfg in README.md #548 (Songmu)


## 0.55.0 (2019-02-13)

* [mongodb] apply for mongodb authenticationDatabase #542 (shibacow)
* [linux] consider the case where the width of the Netid column is only 5 in the output of ss #543 (Songmu)
*  [php-fpm] add option to read status from unix domain socket #541 (lufia)


## 0.54.0 (2019-01-10)

* [redis] Change evicted_keys.Diff to true #539 (lufia)
* [redis] Add evicted_keys metric #537 (lufia)
* Add redash api-key option #535 (kyoshidajp)
* [squid] Add metrics to squid #534 (nabeo)
* [postgres] Enable connection without password #533 (kyoshidajp)


## 0.53.0 (2018-11-12)

* [mysql] Use go-mackerel-plugin instead of go-mackerel-plugin-helper #531 (shibayu36)


## 0.52.0 (2018-10-17)

* Set (default) User-Agent header to HTTP requests #528 (astj)
* Build with Go 1.11 #527 (astj)
* Improve jvm error handling #526 (astj)


## 0.51.1 (2018-08-30)

* [postgres]Ignore error to support Aurora #519 (matsuu)


## 0.51.0 (2018-07-25)

* [mysql] Fix decoding transaction ids from mysql innodb status #521 (itchyny)
* add MSSQL plugin #520 (mattn)


## 0.50.0 (2018-06-20)

* [aws-kinesis-streams] Collect (Write|Read)ProvisionedThroughputExceeded metrics correctly #515 (shibayu36)
* [aws-s3-requests] CloudWatch GetMetricStatics parameters #514 (astj)


## 0.49.0 (2018-05-16)

* [aws-rds]support Aurora PostgreSQL engine #512 (matsuu)
* [aws-rds]fix unit for some metrics #511 (matsuu)
* [aws-rds]add BurstBalance metric #510 (matsuu)
* [linux] fix for collectiong ioDrive(FusionIO) diskstats #509 (hayajo)


## 0.48.0 (2018-04-18)

* [linux] collect disk stats of NVMe devices and ignore virtual/removable devices #506 (hayajo)


## 0.47.0 (2018-04-10)

* [aws-ec2-cpucredit] Add T2 unlimited CPU credit metrics #504 (astj)


## 0.46.0 (2018-03-15)

* [Redis] send uptime #502 (dozen)
* [redis] expired_keys change to diff true #501 (dozen)


## 0.45.0 (2018-03-01)

* [postgres] Add amount of xlog_location change #476 (kizkoh)


## 0.44.0 (2018-02-08)

* [aws-elasticsearch] support metric-{key,label}-prefix #495 (astj)
* [mongodb] Fix warning message on MongoDB 3.4, 3.6 #493 (hayajo)
* Add mackerel-plugin-aws-s3-requests #492 (astj)
* Migrate from `go-mgo/mgo` to `globalsign/mgo` #491 (hayajo)


## 0.43.0 (2018-01-23)

* Setting password via environment variable #488 (hayajo)
* update rpm-v2 task for building Amazon Linux 2 package #485 (hayajo)
* Support BSD #487 (miwarin)
* make `make build` works for some plugins which moved out from this repository #486 (astj)


## 0.42.0 (2018-01-10)

* Move mackerel-plugin-json to other repository #482 (shibayu36)
* Move mackerel-plugin-gearmand #480 (shibayu36)
* Move to mackerelio/mackerel-plugin-gcp-compute-engine #479 (shibayu36)
* [mongodb] fix connections_current metric mongodb-Replica-Set #467 (vfa-cancc)
* [haproxy]support unix domain socket #477 (hbadmin)
* [postgres]state may be null even in old versions #478 (matsuu)
* [uptime] use go-osstat/uptime instead of golib/uptime for getting more accurate uptime #475 (Songmu)
* [mysql] add a hint for -disable_innodb #474 (astj)


## 0.41.1 (2017-12-20)

* [mysql] set Diff: true for some stats which are actually counter values #472 (astj)


## 0.41.0 (2017-12-20)

* [mysql] Fix some InnoDB stats #469 (astj)
* [mysql] Fix message for socket option #468 (utisam)
* MySQL Plugin support Aurora reader node #462 (dozen)


## 0.40.0 (2017-12-12)

* Add h2o to package #464 (astj)
* Redis Plugin supports custom CONFIG command #463 (dozen)
* add mackerel-plugin-h2o #456 (hayajo)
* add defer to closing the response body, and change position it. #461 (qt-luigi)
* add that close the response body #460 (qt-luigi)
* [redis] Add Redis replication delay and lag metrics #455 (kizkoh)


## 0.39.0 (2017-11-28)

* Don't add plugins README which has been moved #453 (astj)
* Improve docker plugin #436 (astj)
* [jvm] Fix remote jvm monitoring #451 (astj)
* Changed README.md of mackerel-plugin-linux #452 (soudai)
* [json] Fix error handling #449 (astj)
* Fix license notice #448 (itchyny)
* [docker] Avoid concurrent map writes by multiple goroutines #446 (astj)
* [aws-ec2-ebs] Do not log "fetched no datapoints" error #445 (astj)
* [kinesis-streams] Use Sum aggregation for Kinesis streams statistics #435 (itchyny)


## 0.38.0 (2017-11-09)

* Improve mackerel-plugin-postgres #437 (astj)
* [docker] Add CPU Percentage metrics #424 (astj)
* [gostats] Use go-mackerel-plugin instead of go-mackerel-plugin-helper #429 (itchyny)
* [mysql]Fix makeBigint calculation #434 (matsuu)
* [cloudfront] add -metric-key-prefix option #433 (fujiwara)


## 0.37.1 (2017-10-26)

* [multicore] Refactor multicore plugin #430 (itchyny)


## 0.37.0 (2017-10-19)

* Implement mackerel-plugin-mcrouter #420 (waniji)
* [uptime] use go-mackerel-plugin instead of using go-mackerel-plugin-helper #428 (Songmu)


## 0.36.0 (2017-10-12)

* Add mackerel-plugin-json #395 (doublemarket)
* [awd-dynamodb] [incompatible] remove `.` from Metrics.Name #423 (astj)
* [unicorn] Support metric-key-prefix #425 (astj)
* [aws-elasticsearch] Improve CloudWatch Statistic type and add some metrics #387 (holidayworking)


## 0.35.0 (2017-10-04)

* [twemproxy] [incompatible] add `-enable-each-server-metrics` option #419 #421 (Songmu)


## 0.34.0 (2017-09-27)

* add mackerel-plugin-flume to package #415 (y-kuno)
* [mysql]add MyISAM related graphs #406 (matsuu)
* add mackerel-plugin-sidekiq to package #417 (syou6162)
* build with Go 1.9 #414 (astj)
* [OpenLDAP] fix get latestCSN #413 (masahide)
* [aws-dynamodb] Add ReadThrottleEvents metric and fill 0 when *ThrottleEvents metrics are not present #409 (astj)


## 0.33.0 (2017-09-20)

* add mackerel-plugin-nvidia-smi to package #411 (syou6162)
* [accesslog] Feature/accesslog/customize parser #410 (karupanerura)
* Fix redundant error by golint in redis.go #408 (shibayu36)
* add flume plugin #396 (y-kuno)
* [mysql]add handler graphs #402 (matsuu)


## 0.32.0 (2017-09-12)

* [memcached] add evicted.reclaimed and evicted.nonzero_evictions #388 (Songmu)
* [mysql]add missed metrics and fix graph definition #390 (matsuu)
* [Redis] fix expired keys #398 (edangelion)
* [accesslog] Fix for scanning long lines #400 (itchyny)


## 0.31.0 (2017-08-30)

* [redis] Change queries metric to diff of "total_commands_processed" #397 (edangelion)
* [aws-dynamodb] Refactor and parallelize CloudWatch request with errgroup #367 (astj)
* [plack] Don't raise errors when parsing JSON fields failed #394 (astj)
* [jmx-jolokia] add value to thread graph #393 (y-kuno)


## 0.30.0 (2017-08-23)

* add mackerel-plugin-openldap to package #391 (astj)
* Add Burst Balance metric for AWS EC2 EBS plugin #384 (ariarijp)
* Add openldap plugin  #374 (masahide)


## 0.29.1 (2017-08-02)

* [solr] Fix a graph definition for Apache Solr's cumulative metric #381 (supercaracal)
* [accesslog] Refine LTSV format detection logic https://github.com/Songmu/axslogparser/pull/8 (Songmu)
* [accesslog] Fix testcase (Percentile logic is Fixed up) #380 (Songmu)


## 0.29.0 (2017-07-26)

* [aws-dynamodb] Add TimeToLiveDeletedItemCount metrics #376 (astj)
* [aws-dynamodb] Adjust options and graph definitions #375 (astj)
* [mysql] Fix graph label prefixes #372 (koooge)


## 0.28.1 (2017-06-28)

* postgres: add metric-key-prefix #363 (edangelion)
* [accesslog] add mackerel-plugin-accesslog #359 (Songmu)
* add mackerel-plugin-aws-dynamodb to package #366 (astj)
* Use mackerelio/golib/logging as logger, not mackerelio/mackerel-agent/logging #365 (astj)
* postgres: collect dbsize only if connectable #361 (mechairoi)
* Support PostgreSQL 9.6 #360 (mechairoi)
* Add sidekiq plugin #354 (littlekbt)


## 0.28.0 (2017-06-14)

* Add aws-dynamodb plugin #349 (astj)
* Implemented mackerel-plugin-redash #355 (yoheimuta)
* Add mackerel-plugin-solr to package #356 (astj)
* Add test cases and fix issues for apache solr #341 (supercaracal)


## 0.27.2 (2017-06-07)

* disable diff on php-opcache.cache_size because they are gauge value #352 (matsuu)
* build with Go 1.8 #350 (Songmu)
* v2 packages (rpm and deb) #348 (Songmu)
* [aws-rds] Fix "Latency" metric label #347 (astj)
* Add AWS Kinesis Firehose Plugin #333 (holidayworking)
* Fixed mackerel-plugin-nginx/README.md #345 (kakakakakku)


## 0.27.1 (2017-05-09)

* [php-fpm] Implement PluginWithPrefix interfarce #338 (astj)
* Use SetTempfileByBasename to support MACKEREL_PLUGIN_WORKDIR #339 (astj)


## 0.27.0 (2017-04-27)

* Add uWSGI vassal plugin #335 (kizkoh)
* add mackerel-plugin-uwsgi-vassal to package #336 (astj)


## 0.26.0 (2017-04-19)

* Add AWS Rekognition Plugin #322 (holidayworking)
* Add aws-kinesis-streams plugin #326 (astj)
* Add AWS Lambda plugin #327 (astj)
* [redis] fix metrics lable #329 (y-kuno)
* Add aws-lambda and aws-kinesis-streams to package #330 (astj)
* Support twemproxy v0.3, Add total_server_error #332 (masahide)


## 0.25.6 (2017-04-06)

* Cross compile by go's native cross build, not by gox #321 (astj)
* fix a label of gostats plugin #323 (itchyny)


## 0.25.5 (2017-03-22)

* add `mackerel-plugin` command #315 (Songmu)
* Add AWS WAF Plugin #316 (holidayworking)
* use new bot token #318 (daiksy)
* use new bot token #319 (daiksy)


## 0.25.4 (2017-02-22)

* Improve gce plugin #313 (astj)


## 0.25.3 (2017-02-16)

* Feature/gcp compute engine #304 (littlekbt)
* [aws-rds] Make it possible to get metrics from Aurora. #307 (TakashiKaga)
* [multicore]fix tempfile path #311 (daiksy)


## 0.25.2 (2017-02-08)

* [aws-rds] fix metric name #306 (TakashiKaga)
* [aws-ses] ses.stats is unit type #308 (holidayworking)
* [aws-cloudfront] Fix regression #295 #309 (astj)


## 0.25.1 (2017-01-25)

* Make more plugins to support MACKEREL_PLUGIN_WORKDIR #301 (astj)
* [jvm] Fix the label and scale #302 (itchyny)
* [aws-rds] Support Aurora metrics and refactoring #303 (sioncojp)


## 0.25.0 (2017-01-04)

* Change directory structure convention of each plugin #289 (Songmu)
* [apache2] fix typo in graphdef #291 (astj)
* [apache2] Change metric name not to end with dot #293 (astj)
* add mackerel-plugin-windows-server-sessions #294 (daiksy)
* migrate from goamz to aws-sdk-go #295 (astj)
* [docker] Add timeout for API request #296 (astj)


## 0.24.0 (2016-11-29)

* Implement mackerel-plugin-aws-ec2 #248 (yyoshiki41)
* [postgres] support Pg9.1 #274 (Songmu)
* Add new nvidia-smi plugin #280 (ksauzz)
* [jvm] Add notice about user to README #281 (astj)
* Implement mackerel-plugin-twemproxy #283 (yoheimuta)
* fix cloudwatch dimensions for elb #284 (ki38sato)
* Change error strings to pass current golint #285 (astj)
* Add mackerel-plugin-twemproxy to package #286 (stefafafan)


## 0.23.1 (2016-10-27)

* [redis] Fix a bug to fetch no metrics of keys and expired #272 (yoheimuta)
* fix: "open file descriptors" property in elasticsearch  #273 (kamijin-fanta)
* [memcached] Supported memcached curr_items metric #275 (kakakakakku)
* [memcached] support new_items metrics #276 (Songmu)
* [redis] s/memoty/memory/ #277 (astj)


## 0.23.0 (2016-10-18)

* mackerel-plugin-linux: Allow to select multiple (but not all) sets of metrics #243 (astj)
* Fixed flag comment of mackerel-plugin-fluentd #257 (kakakakakku)
* Fix postgres.iotime.{blk_read_time,blk_write_time} #259 (mechairoi)
* [Plack] Adopt Plack::Middleware::ServerStatus::Lite 0.35's response #261 (astj)
* build with Go 1.7 #262 (astj)
* Add much graphs/metrics to mackerel-plugin-mysql #264 (netmarkjp)
* [apache2] Support -metric-key-prefix option and get rid of default Tempfile specification #265 (astj)
* [aws-rds] add `-engine` option #266 (Songmu)
* [elasticsearch] Add open_file_descriptors metric in elasticsearch plugin #267 (kamijin-fanta)
* Make *some* plugins to support MACKEREL_PLUGIN_WORKDIR #268 (astj)
* [redis] deal with MACKEREL_PLUGIN_WORKDIR #269 (astj)


## 0.22.1 (2016-09-06)

* Fixed README.md #253 (kakakakakku)
* [memcached] Support -metric-key-prefix option  #255 (astj)


## 0.22.0 (2016-07-14)

* add multicore plugin #246 (daiksy)
* add mackerel-plugin-multicore into package (daiksy)


## 0.21.2 (2016-07-07)

* Fix help message #244 (ariarijp)
* [apache2] update README.md. fix mod_status configuration #245 (Songmu)
* Add some plugins to README #247 (ariarijp)
* follow urfave/cli #249 (Songmu)
* [mysql] support -metric-key-prefix option #250 (Songmu)


## 0.21.1 (2016-06-28)

* build with go 1.6.2 #241 (Songmu)


## 0.21.0 (2016-06-23)

* Add PHP-FPM plugin #226 (ariarijp)
* Support password authentication of Redis #232 (hico-horiuchi)
* Add an option to specify type and id pattern to fluentd plugin #233 (waniji)
* Fix bug:aws-ses #234 (tjinjin)
* fix help link #235 (daiksy)
* xentop: get CPU %, not CPU time/min #236 (hagihala)
* add mackerel-plugin-php-fpm into package #238 (Songmu)


## 0.20.2 (2016-06-09)

* aws-ec2-ebs: Use wildcard in the graph definitions #230 (itchyny)


## 0.20.1 (2016-05-25)

* change signatures of doMain to follow recent codegangsta/cli #227 (Songmu)
* fix README.md of mackerel-plugin-jvm #228 (azusa)


## 0.20.0 (2016-05-10)

* [docker] use goroutine for fetching metrics via API #220 (stanaka)
* add graphite and proc-fd into package #223 (Songmu)


## 0.19.4 (2016-04-20)

* Add mackerel-plugin-graphite (#216) (taku-k)
* Add mackerel plugin proc fd (#207) (taku-k)
* Do not send fluentd metrics of other than the output plugin (#213) (waniji)

## 0.19.3 (2016-04-14)

* [redis] skip to calculate capacity when CONFIG command failed (#214) (Songmu)
* Revert "Revert "use /usr/bin/mackerel-plugin-*"" #208 (Songmu)
* fix: redis plugin panics when redis-server is not installed. #209 (stanaka)
* fix: rpm should not include dir #210 (stanaka)
* [nginx] fix typo #215 (y-kuno)
* Refactoring the release process #217 (stanaka)
