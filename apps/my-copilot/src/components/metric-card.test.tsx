import { render, screen } from "@testing-library/react";
import MetricCard from "./metric-card";

describe("MetricCard", () => {
  const defaultProps = {
    value: 1234,
    label: "Aktive brukere",
    helpText: "Antall brukere som har brukt Copilot siste 30 dager",
    helpTitle: "Aktive brukere",
  };

  it("renders label and value", () => {
    render(<MetricCard {...defaultProps} />);

    expect(screen.getByText("Aktive brukere")).toBeInTheDocument();
    expect(screen.getByText(1234)).toBeInTheDocument();
  });

  it("renders subtitle when provided", () => {
    render(<MetricCard {...defaultProps} subtitle="+12% fra forrige måned" />);

    expect(screen.getByText("+12% fra forrige måned")).toBeInTheDocument();
  });

  it("does not render subtitle when not provided", () => {
    render(<MetricCard {...defaultProps} />);

    expect(screen.queryByText("+12% fra forrige måned")).not.toBeInTheDocument();
  });

  it("uses smaller heading for long text values", () => {
    render(<MetricCard {...defaultProps} value="1 234 567" />);

    const heading = screen.getByRole("heading", { level: 2 });
    expect(heading).toHaveTextContent("1 234 567");
    expect(heading).toHaveClass("aksel-heading--medium");
  });

  it("uses xlarge heading for short values", () => {
    render(<MetricCard {...defaultProps} value={42} />);

    const heading = screen.getByRole("heading", { level: 2 });
    expect(heading).toHaveClass("aksel-heading--xlarge");
  });
});
