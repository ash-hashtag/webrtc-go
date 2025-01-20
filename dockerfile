FROM coturn/coturn

# Add two sets of credentials
RUN echo "username1:password1" >> /etc/turnuserdb.conf
RUN echo "username0:password0" >> /etc/turnuserdb.conf

# Expose the required ports
EXPOSE 3478
EXPOSE 3478/udp
EXPOSE 5349
EXPOSE 5349/udp

CMD ["turnserver"]
