FROM rust:1.84.1 AS build

ENV RUST_LOG="info"
ENV TELOXIDE_TOKEN="REDACTED"

WORKDIR /opt

COPY ./bot ./

RUN ls -la

RUN cargo build --release --bin main

FROM alpine:latest AS runner

COPY --from=build /opt/target/release/main /

RUN chmod +x /main

CMD ["/main"]

