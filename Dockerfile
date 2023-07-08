FROM scratch
ADD conf /conf
ADD web/build /build
ADD k8swt /
CMD ["/k8swt"]