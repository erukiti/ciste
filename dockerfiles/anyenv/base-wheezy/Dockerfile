#docker build -t erukiti/anyenv:base-wheezy .

FROM buildpack-deps:wheezy

WORKDIR /root
RUN git clone https://github.com/riywo/anyenv /root/.anyenv
RUN echo 'export PATH="$HOME/.anyenv/bin:$PATH"' >> ~/.bash_profile
RUN echo 'eval "$(anyenv init -)"' >> ~/.bash_profile
ENV PATH="/root/.anyenv/bin:$PATH"
RUN anyenv init -
RUN bash -c eval "$(anyenv init -)"
