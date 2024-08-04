#######################################################################################
#                                                                                     #
# ######   #######      #     #  #######  #######      #######  ######   ###  ####### #
# #     #  #     #      ##    #  #     #     #         #        #     #   #      #    #
# #     #  #     #      # #   #  #     #     #         #        #     #   #      #    #
# #     #  #     #      #  #  #  #     #     #         #####    #     #   #      #    #
# #     #  #     #      #   # #  #     #     #         #        #     #   #      #    #
# #     #  #     #      #    ##  #     #     #         #        #     #   #      #    #
# ######   #######      #     #  #######     #         #######  ######   ###     #    #
#                                                                                     #
#######################################################################################
#
# This is the catchall makefile. Its purpose is to invoke the default rules if
# you do not have a given rule defined in your own Makefile.
#
# You should not have any reason to edit this makefile. If you want to prevent
# it from triggering and/or invoking default rules, you can just define your own
# rules for whatever, including empty rules.
#
# If you want something in here changed for whatever reason, ask rops. If you
# start making changes to this file, rops will be disappointed with you.

# NOTE(dmr, 2021-04-27): in the repo we symlink .default.mk to this file, so you
# can more easily test & play with it.
DEFAULTS_MAKEFILE ?= .default.mk

# official example; required for old-make compat?
#
# NOTE(dmr, 2021-04-27): additional context: we want to make sure that this
# stuff runs on as wide a set of make versions as possible. For example, MacOS
# ships with make version 3.81 or something, which I believe was the version
# released around the Fall of Rome. Nonetheless, we would like this to all work
# just fine on the dev machines of MacOS users who have not bothered to install
# less archaic versions, as well as in linux CI envs.
#
# The commented example below is, IIRC, targeted at an older version of MacOS.
# I'm keeping it around so its handy in case the uncommented example, which has
# not been extensively tested on pre-Cambrian versions of make, starts to cause
# issues. Then we'll be more likely to be able to resolve it faster.
#
# %: force
# 	@$(MAKE) -f $(DEFAULTS_MAKEFILE) $@
#
# force: ;

%:
	@$(MAKE) --no-print-directory -f $(DEFAULTS_MAKEFILE) $@
