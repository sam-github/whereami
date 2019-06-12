# whereami?

Given a directory of some images, that contains a sub-directory that contains
some images, create a command line utility that reads the EXIF data from the
images and writes the image path, latitude and longitude to file as a CSV. Use
Go routines and channels to do the reads concurrently if possible. For extra
credit, provide an option to write to HTML as well.
