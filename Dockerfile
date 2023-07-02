FROM scratch
ADD conf /
ADD web/dist /
ADD k8swt /
CMD ["/k8swt"]