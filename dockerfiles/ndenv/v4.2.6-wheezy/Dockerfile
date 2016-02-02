#docker build -t ndenv:4.2.6-wheezy .

FROM anyenv:base-wheezy

RUN anyenv install ndenv
RUN bash -l -c "ndenv install v4.2.6"
RUN bash -l -c "ndenv global v4.2.6"
