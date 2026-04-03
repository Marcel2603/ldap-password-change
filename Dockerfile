ARG TARGETPLATFORM

FROM scratch
ARG TARGETPLATFORM
ENTRYPOINT ["/usr/bin/ldap-password-change"]
EXPOSE 8080
COPY static /usr/bin/static
COPY $TARGETPLATFORM/ldap-password-change /usr/bin/
