# Makefile — PDP‑8 Lunar Lander (FOCAL‑69) demo with Fedora’s simh‑pdp8
# Prerequisites (Fedora 42):
#   sudo dnf install simh curl make tar gzip
# Usage:
#   make        # fetch everything + launch SIMH ready for paste‑in of Lunar code
#   make clean  # tidy downloaded artefacts

# ───────────────────────────── Config ─────────────────────────────
PDP8            ?= simh-pdp8                                   # simulator binary

# ░░ FOCAL‑69 interpreter ░░
FOCAL_TAR_URL  := https://downloads.sourceforge.net/project/simh/Software%20Kits/FOCAL69%20for%20the%20PDP-8./foclswre.tar.z
FOCAL_TAR      := foclswre.tar.Z                              # Unix compress format
FOCAL_BIN      := focal69.bin                                 # extracted binary paper‑tape image

# ░░ Lunar Lander source in FOCAL (transcribed plain text) ░░
LUNAR_SRC_URL  := https://www.cs.brandeis.edu/~storer/LunarLander/LunarLander/LunarLanderListingText.txt
LUNAR_SRC      := lunar-lander.foc

INI            := pdp8.ini

# ─────────────────────────── Targets ──────────────────────────────
.PHONY: all run clean untar

all: run

$(FOCAL_TAR):
	curl -L -o $@ $(FOCAL_TAR_URL)

# Extract only focal69.bin from the .tar.Z (uses POSIX tools)
#
# broken
# $(FOCAL_BIN): $(FOCAL_TAR)
#	@echo "# extracting $@ from $< …"
#	uncompress -c $(FOCAL_TAR) | tar -xOf - ./focal69.bin > $@

$(LUNAR_SRC):
	curl -L -o $@ $(LUNAR_SRC_URL)

$(INI): $(FOCAL_BIN)
	@echo "; auto‑generated pdp8.ini"              >  $@
	@echo "SET CPU 4K"                              >> $@
	@echo "SET CPU NOEAE"                           >> $@
	@echo "LOAD $(FOCAL_BIN)"                       >> $@
	@echo "RUN 200"                                >> $@

# the simh based simulator is not really working in this interactive game
# simh: $(FOCAL_BIN) $(LUNAR_SRC) $(INI)
simh:
	@echo "\n*** Launching SIMH PDP‑8 with FOCAL‑69 …\n"
	$(PDP8) -i $(INI)
	@echo "\nPaste the contents of $(LUNAR_SRC) at the '*' prompt, then type GO.\n"

clean:
	rm -f $(FOCAL_TAR) $(FOCAL_BIN) $(LUNAR_SRC) $(INI)

.PHONY: run
run:
	retrofocal lunar-lander.fc
