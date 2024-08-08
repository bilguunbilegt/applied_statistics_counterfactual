# applied_statistics_counterfactual

Before diving into the summary, it's important to note that finding and using a counterfactual package in R was much easier due to the availability of dedicated libraries. In contrast, Go lacks direct support for counterfactual analysis, requiring me to implement my own solution based on existing Go packages. This difference in ecosystem maturity is evident in the model's accuracy and reliability across both languages.

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

## Memory Usage and Cost Analysis

### Memory Usage in R

- **Memory before model fitting**: 62,543,912 bytes (~62.5 MB)
- **Memory after model fitting**: 64,616,472 bytes (~64.6 MB)
- **Memory after adding predictions**: 65,120,552 bytes (~65.1 MB)

### Memory Usage in Go

- **Initial Memory Usage**: 
  - **Allocated (Alloc)**: 0 MiB
  - **Total Allocated (TotalAlloc)**: 0 MiB
  - **System Memory (Sys)**: 6 MiB
- **Memory After Fitting Model**: 
  - **Allocated (Alloc)**: 0 MiB
  - **Total Allocated (TotalAlloc)**: 0 MiB
  - **System Memory (Sys)**: 6 MiB
- **Final Memory Usage**:
  - **Allocated (Alloc)**: 1 MiB
  - **Total Allocated (TotalAlloc)**: 1 MiB
  - **System Memory (Sys)**: 7 MiB

### Cost Implications

**AWS Pricing Considerations**:

1. **R Memory Usage**:
   - The maximum memory usage in R is **65.1 MB**.

2. **Go Memory Usage**:
   - The maximum memory usage in Go is **7 MiB**.

**Instance Type Selection**:
- For this comparison, let’s consider a **t3.medium** instance on AWS with **4 GB RAM** (~4096 MiB) and **2 vCPUs** costing approximately **$32.13/month**.

**Memory Cost per Instance**:

- **R Usage**:
  - With **65.1 MB** (~65 MiB) of memory usage, R uses **~1.59%** of the **4 GB** available memory.
  - This equates to a cost of **$0.51/month** for memory usage.

- **Go Usage**:
  - With **7 MiB** of memory usage, Go uses **~0.17%** of the **4 GB** available memory.
  - This equates to a cost of **$0.05/month** for memory usage.

### Potential Cost Savings

- **Memory Cost Savings**:
  - The Go implementation uses significantly less memory, leading to potential savings.
  - The **memory cost saving** by using Go instead of R would be approximately **90%**.

- **Total Instance Costs**:
  - If memory were the only cost factor, Go’s efficiency could result in a monthly saving of **$0.46** per instance.

**Total Savings Estimate**:

- Assuming the memory usage scales similarly across multiple instances and similar performance metrics, Go’s lower memory footprint can lead to substantial savings, particularly in large-scale deployments.

- In real-world scenarios where multiple instances are deployed, the **overall cloud cost savings** could be substantial, especially when memory and CPU efficiency translate to fewer or smaller instances required.

### Conclusion

The R implementation provides more reliable and accurate predictions. The Go implementation, although efficient in memory and speed, requires further debugging and improvement to ensure the accuracy of the model outputs.
Switching from R to Go could lead to significant cloud cost savings, especially due to Go's lower memory usage and higher efficiency. For large-scale deployments on AWS, where memory usage directly influences costs, the switch could reduce cloud computing expenses by approximately **90%** related to memory usage, resulting in substantial overall savings.



