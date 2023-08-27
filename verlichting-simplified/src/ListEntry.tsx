import { PropsWithChildren } from "react";
import styled from "styled-components";

type Props = {
    isActive: boolean;
    isLoading: boolean;
    onClick?: () => void;
};

export const ListEntry = ({ isActive, isLoading, onClick, children }: PropsWithChildren<Props>) => {
    return <Li isActive={isActive} isLoading={isLoading}><Link onClick={onClick} isActive={isActive} isLoading={isLoading}>{children}</Link></Li>
};

const Link = styled.a<Props>`
    display: flex; /* To make the element fill its container */

    padding: 10px;

    line-height: 20px;

    text-decoration: none;

    cursor: pointer;

    background-color: ${(props) => (props.isActive ? "#3333cc" : props.isLoading ? "#9999ff" : "white")};
    color: ${(props) => (props.isActive || props.isLoading ? "white" : "black")};

    &:hover {
        color: ${(props) => (props.isActive ? "white" : "black")};
        background-color: ${(props) => (props.isActive ? "#3333cc" : props.isLoading ? "#7777ee" : "#bbbbff")};
    }
`;

const Li = styled.li<Props>`
    list-style: none;
    text-align: left;

    border: 1px solid black;
    border-radius: 5px;
    margin: 6px 0px;
`;
