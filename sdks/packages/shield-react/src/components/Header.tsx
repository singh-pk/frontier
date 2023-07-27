import { Flex, Text } from "@raystack/apsara";
import React, { ComponentPropsWithRef } from "react";
import { useFrontier } from "../contexts/FrontierContext";

const styles = {
  container: {
    fontSize: "12px",
    minWidth: "220px",
    maxWidth: "100%",
    color: "var(--foreground-base)",

    display: "flex",
    flexDirection: "column",
    alignItems: "center",
    gap: "32px",
  },
  logoContainer: {},
  titleContainer: {
    fontWeight: "400",
  },
};

const defaultLogo = (
  <img
    src="data:image/png;base64,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"
    style={{ borderRadius: "var(--pd-8)", width: "80px", height: "80px" }}
  ></img>
);

type HeaderProps = ComponentPropsWithRef<"div"> & {
  title?: string;
  logo?: React.ReactNode;
};

export const Header = ({ title, logo }: HeaderProps) => {
  const { config } = useFrontier();

  return (
    <Flex
      style={{
        ...styles.container,
        flexDirection: "column",
      }}
    >
      <div style={styles.logoContainer}>{logo ? logo : defaultLogo}</div>
      <div style={styles.titleContainer}>
        <Text size={9}>{title}</Text>
      </div>
    </Flex>
  );
};
