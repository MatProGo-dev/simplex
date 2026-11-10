# `linprog`

This package provides a Go dual revised-simplex solver for continuous linear
programs represented by `problem.OptimizationProblem`. It converts supported
models to standard form, initializes a dual-feasible basis through an auxiliary
Phase-I feasibility problem, and applies dual-simplex pivots until primal
feasibility and optimality are established.

The method is inspired by the dual revised-simplex implementation referenced
by SciPy's `linprog(method="highs-ds")`:

- Huangfu, Q. and Hall, J. A. J., “Parallelizing the dual revised simplex
  method,” *Mathematical Programming Computation*, 10 (1), 119–142 (2018).
  [doi:10.1007/s12532-017-0130-5](https://doi.org/10.1007/s12532-017-0130-5)
- [HiGHS high-performance linear optimization software](https://highs.dev/)

This is an independent, serial Go implementation, not a port or wrapper of
HiGHS. It does not implement SciPy's `highs-ipm`, automatic method selection,
presolve, mixed-integer solving, Python options, or SciPy's detailed result
and sensitivity fields. The algorithm currently uses fixed feasibility
tolerances and does not provide dual certificates or sensitivity information.
