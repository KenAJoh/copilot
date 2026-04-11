import { render, screen, fireEvent } from "@testing-library/react";
import { DomainCard } from "./domain-card";

describe("DomainCard", () => {
  const defaultProps = {
    domain: "platform" as const,
    count: 5,
    selected: false,
    onClick: vi.fn(),
  };

  beforeEach(() => {
    vi.clearAllMocks();
  });

  it("renders domain label and count", () => {
    render(<DomainCard {...defaultProps} />);

    expect(screen.getByRole("heading", { level: 3 })).toBeInTheDocument();
    expect(screen.getByText("5")).toBeInTheDocument();
  });

  it("calls onClick with domain when clicked", () => {
    render(<DomainCard {...defaultProps} />);

    fireEvent.click(screen.getByRole("button"));
    expect(defaultProps.onClick).toHaveBeenCalledWith("platform");
  });

  it("applies selected styling when selected", () => {
    const { container } = render(<DomainCard {...defaultProps} selected={true} />);

    const button = container.querySelector("button");
    expect(button?.className).toContain("border-blue-500");
  });

  it("does not apply selected styling when not selected", () => {
    const { container } = render(<DomainCard {...defaultProps} selected={false} />);

    const button = container.querySelector("button");
    expect(button?.className).toContain("border-transparent");
  });

  it("renders different domains", () => {
    render(<DomainCard {...defaultProps} domain="frontend" count={12} />);

    expect(screen.getByText("12")).toBeInTheDocument();
  });
});
