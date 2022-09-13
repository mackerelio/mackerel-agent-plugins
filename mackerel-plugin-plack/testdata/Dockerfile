FROM perl:5.36

COPY . /usr/src/myapp
WORKDIR /usr/src/myapp

RUN cpm install -v

EXPOSE 5000

ENV PERL5LIB /usr/src/myapp/local/lib/perl5
ENV PATH $PATH:/usr/src/myapp/local/bin/

CMD [ "plackup", "app.psgi" ]

