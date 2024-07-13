export type Domain = {
  ID: number;
  Type: "A" | "CNAME" | "MX" | "";
  Domain: string;
  Record: string;
  Status: 0 | 1;
  Result?: ParseResult;
};

export type ParseResult = {
  ISP?: string;
  IP?: string;
  City?: string;
  ExpireAt?: string;
  Remark?: string;
};

export const Domains: Domain[] = [
  {
    ID: 1,
    Type: "A",
    Domain: "www.example.com",
    Record: "172.168.100.100",
    Status: 1,
    Result: {
      ISP: "Example ISP",
      IP: "172.168.100.100",
      City: "Example City",
      ExpireAt: "2021-12-31",
      Remark: "Example Remark",
    },
  },
  {
    ID: 2,
    Type: "A",
    Domain: "www.example.com",
    Record: "172.168.100.100",
    Status: 1,
    Result: {
      ISP: "Example ISP",
      IP: "172.168.100.100",
      City: "Example City",
      ExpireAt: "2021-12-31",
      Remark: "Example Remark",
    },
  },
  {
    ID: 3,
    Type: "A",
    Domain: "www.example.com",
    Record: "172.168.100.100",
    Status: 1,
    Result: {
      ISP: "Example ISP",
      IP: "172.168.100.100",
      City: "Example City",
      ExpireAt: "2021-12-31",
      Remark: "Example Remark",
    },
  },
];
