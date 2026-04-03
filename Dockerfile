ARG TARGETPLATFORM

FROM scratch
ARG TARGETPLATFORM
WORKDIR /opt/ldap-password-change
ENTRYPOINT ["/opt/ldap-password-change/service"]
EXPOSE 8080
COPY ldap-password-change /opt/ldap-password-change/service
