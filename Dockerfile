FROM quay.io/eris/base
MAINTAINER Eris Industries <support@erisindustries.com>

ENV NAME         toadserver
ENV REPO 	 eris-ltd/$NAME
ENV BRANCH       master
ENV BINARY_PATH  $NAME
ENV CLONE_PATH   $GOPATH/src/github.com/$REPO
ENV INSTALL_PATH $INSTALL_BASE/$NAME

RUN mkdir -p $CLONE_PATH

#for local buildz
#COPY . $CLONE_PATH
#WORKDIR $CLONE_PATH

RUN git clone -q https://github.com/$REPO $CLONE_PATH
RUN git checkout -q $BRANCH
RUN go install

USER $USER
WORKDIR $ERIS

VOLUME $ERIS
EXPOSE 11113
CMD ["toadserver"]
