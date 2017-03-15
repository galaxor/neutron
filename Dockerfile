FROM scratch
WORKDIR /
ADD neutron-docker /neutron
CMD ["/neutron"]
EXPOSE 4000
