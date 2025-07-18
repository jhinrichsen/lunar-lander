# Makefile – clean Lunar Lander scans for OCR and extract text
# Tested on Fedora 42 (DNF5)
# Usage examples:
#   make            # install deps + clean + OCR all default images
#   make ocr        # just run the image pipeline
#   make clean      # remove output directory

SRC      = LunarLanderListing.jpg \
          LunarLanderSampleOutputPage1.jpg \
          LunarLanderSampleOutputPage2.jpg
OUTDIR   = ocr_ready
#–– Tunables ––––––––––––––––––––––––––––––––––––––––––––––––––––
TARGET_DPI     ?= 600
DESKEW         ?= 40%%
UNSHARP        ?= 0x0.8+0.8+0.02
BIN_THRESH     ?= 55%%
CONTRAST       ?= 1%%x1%%
#———————————————————————————————————————————————————————————————

.PHONY: all
all: ocr

.PHONY: deps
deps:
	@echo "Installing OCR dependencies (ImageMagick, unpaper, tesseract)…"
	sudo dnf install -y ImageMagick unpaper tesseract

$(OUTDIR):
	@mkdir -p $(OUTDIR)

.PHONY: ocr
ocr: $(OUTDIR) $(SRC)
	@for img in $(SRC); do \
	  base=$$(basename $$img .jpg); \
	  tmp1=$(OUTDIR)/$${base}_gray.tif; \
	  tmp2=$(OUTDIR)/$${base}_deskew.tif; \
	  tmp3=$(OUTDIR)/$${base}_clean.tif; \
	  final=$(OUTDIR)/$${base}_final.tif; \
	  echo "▶  Cleaning $$img"; \
	  magick $$img -units PixelsPerInch -density $(TARGET_DPI) -colorspace Gray -auto-level $$tmp1; \
	  magick $$tmp1 -deskew $(DESKEW) +repage $$tmp2; \
	  magick $$tmp2 \
	    -morphology Convolve Gaussian:2,5 \
	    -sharpen $(UNSHARP) \
	    -contrast-stretch $(CONTRAST) $$tmp3; \
	  magick $$tmp3 -threshold $(BIN_THRESH) $$final; \
	  echo "    → OCR $$final → $${base}.txt"; \
	  tesseract $$final $(OUTDIR)/$${base} --psm 6 txt >/dev/null 2>&1; \
	  rm -f $$tmp1 $$tmp2 $$tmp3; \
	done
	@echo "✓  OCR results written to $(OUTDIR)/*.txt"

.PHONY: clean
clean:
	rm -rf $(OUTDIR)
