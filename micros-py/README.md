
//  python -m grpc_tools.protoc -I api/pipeline/v1 --python_out=api/pipeline/v1 --grpc_python_out=api/pipeline/v1 --pyi_out=api/pipeline/v1 api/pipeline/v1/pipeline.proto
// python -m grpc_tools.protoc -I api/dwdata/v2 --python_out=api/dwdata/v2 --grpc_python_out=api/dwdata/v2 --pyi_out=api/dwdata/v2 api/dwdata/v2/dw_data.proto

# contentV2变化
### 新增字段:
        CustomerDetail25,
        CustomerDetail26,
        CustomerDetail27,
        SellingSTA36,
        SupplierRanking25,
        SupplierRanking26,
        SupplierRanking27,
        PurchaseSTA36。
        YdcgqkInfo[
            GNCGAVG4
            JKCGAVG4
            HJAVG4
        ]
        NdcgqkInfo[
            GNCG4
            HJ4
            TITLE4
        ]
        YdxsqkInfo[
            SBAVG3
            CYL3
            FPXXAVG3
        ]

### 减少字段：
    customerDetail12
    supplierRanking12
    


