FROM scratch
ADD conf /conf
ADD web/dist /dist
ADD k8swt /
CMD ["/k8swt"]