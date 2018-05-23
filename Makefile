vendor: Gopkg.lock
	dep ensure -v -vendor-only

Gopkg.lock: Gopkg.toml
	dep ensure -v -no-vendor

.PHONY: dep-update
dep-update: clean-vendor
	dep ensure -update

# clean-vendor has not been added to vendor dependencies as dep is able to check out appropriate
# packages on versions set in the Gopkg.lock file. Removing vendor would force dep to re-download
# all the packages instead of only the missing ones. Global prune is turned off (see Gopkg.toml
# for explanation) thus vendor recipe will leave unused packages in the vendor/ directory. If that
# bothers you, run sequentially clean-vendor and vendor recipes.
.PHONY: clean-vendor
clean-vendor:
	rm -rf vendor
