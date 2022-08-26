FROM alpine:3.16

RUN apk add --update net-snmp net-snmp-tools

EXPOSE 161

CMD ["snmpd", "-V", "-f", "-c", "/etc/snmp/snmpd.conf"]

