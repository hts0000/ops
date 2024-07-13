import React from "react";

const getDomainData = async () => {
  const resp = await fetch("http://localhost:3000/domain-check/api");
  if (!resp.ok) {
    // This will activate the closest `error.js` Error Boundary
    throw new Error("Failed to fetch domain data");
  }

  return resp.json();
};

const DomainTable = async () => {
  const data = await getDomainData();
  console.log(data);
  return <div>DomainTable</div>;
};

export default DomainTable;
