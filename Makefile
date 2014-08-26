all: rpm

rpm:
	rpmbuild --define "_sourcedir `pwd`" -ba packaging/rpm/mackerel-agent-plugins.spec 

