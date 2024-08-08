# applied_statistics_counterfactual

# Summary

This document presents the results of a comparison between R and Go implementations for predicting counterfactual outcomes in a regression model.

## Results Overview

### R Implementation
- **Mean Observed Outcome**: 47.96
- **Mean Counterfactual Outcome (Treatment = 1)**: 47.95
- **Mean Counterfactual Outcome (Treatment = 0)**: 47.98

**Model Coefficients**:
- Intercept: 50.83140
- Treatment: -0.02403 (not statistically significant)
- Covariate: -0.57336 (highly statistically significant)

**Model Performance**:
- Residual Standard Error: 13.87
- R-squared: 0.006681 (Adjusted: 0.006482)
- F-statistic: 33.62, p-value: 2.81e-15

### Go Implementation
- **Mean Observed Outcome**: 47.71
- **Mean Counterfactual Outcome (Treatment = 1)**: 47.28
- **Mean Counterfactual Outcome (Treatment = 0)**: 48.18

**Model Coefficients**:
- Intercept: 50.6319
- Treatment: -0.8976
- Covariate: -0.4847

**Model Performance**:
- Residual Standard Error: 246.70
- R-squared: -352.3825 (Adjusted: -353.4458)
- F-statistic: -0.9972, p-value: 1.0000

## Key Differences

1. **Model Fit**: The R model has a much better fit, indicated by a lower residual standard error and a positive R-squared value, whereas the Go model has a much higher residual standard error and negative R-squared values.

2. **Coefficient Estimates**: While the coefficients in the R model show statistical significance for the covariate, the Go model's coefficients are associated with p-values of 0.0000, which may indicate a calculation error or improper model specification.

3. **Memory and Performance**: The Go implementation demonstrates efficient memory usage and fast execution time, but the accuracy of the results is significantly lower compared to the R implementation.

## Conclusion

The R implementation provides more reliable and accurate predictions. The Go implementation, although efficient in memory and speed, requires further debugging and improvement to ensure the accuracy of the model outputs.
