#
# jeeves
#

# can't use 'scratch' because of the non-portable ref to 'user'
#FROM scratch
FROM debian:latest

EXPOSE 8080

ENTRYPOINT [ \
  "/jeeves", \
  "-port", ":8080" ]

ADD jeeves /jeeves
