FROM golang

RUN mkdir /app

COPY . /app

# change run mode so we don't listen on localhost et al.
RUN sed -i 's/runmode = dev/runmode = prod/' /app/conf/app.conf

WORKDIR /app

RUN go install

EXPOSE 9000

CMD demoweb