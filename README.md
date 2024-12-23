# Golang Fluid Simulation

Sandbox project meant for exploring, learning, and understanding various
aspects of computational fluid dynamics.

Inspired by the following books, papers, etc:
* [Real-Time Fluid Dynamics for Games by Jos Stam](https://www.dgp.toronto.edu/public_user/stam/reality/Research/pdf/GDC03.pdf)
* [Fluid Simulation for Dummies by Mike Ash](https://mikeash.com/pyblog/fluid-simulation-for-dummies.html)
* [Fluid Flow for the Rest of Us: Tutorial of the Marker and Cell Method in Computer Graphics](https://cg.informatik.uni-freiburg.de/intern/seminar/gridFluids_fluid_flow_for_the_rest_of_us.pdf)
* Fluid Simulation for Computer Graphics by Robert Bridson
* An introducion to Computational Fluid Dynamics: The Finite Volume Method by H.K. Versteed & W. Malalasekera

The current iteration of this project is a 2D implementation of an
incompressible fluid system using the staggered Marker-and-Cell Grid data
structure.

## Screenshots 
![01](/screenshots/DensityField.png "01")
![02](/screenshots/VelocityField.png "02")

## Getting Started

### Prerequisites

Golang 1.22
```sh
# Download from https://go.dev or install via homebrew
brew install go
```

### Building and Running

Standard build:
```sh
# If you have make installed, then build and run the app via:
make run

# To build and run without make:
go build -o ./bin/app -v
./bin/app

# Or run without build:
go run .
```

Additional build options:
```sh
# Clean the workspace, deleting binaries, intermediate files, etc:
make clean

# Run with profiling enabled:
make runprof
# View profiler results:
make pprof

# Profiler results require graphviz to be installed:
brew install graphviz
```

## Built With

* [Golang](https://go.dev)
* [Raylib](https://www.raylib.com)

