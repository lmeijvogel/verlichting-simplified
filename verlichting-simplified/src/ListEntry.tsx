import styled from "styled-components";

type Props = {
    isActive: boolean;
    isLoading: boolean;
};

export const ListEntry = styled.li<Props>`
    list-style: none;
    text-align: left;

    line-height: 20px;
    border: 1px solid black;
    padding: 10px;
    border-radius: 5px;
    margin: 6px 0px;
    cursor: pointer;
    color: ${(props) => (props.isActive || props.isLoading ? "white" : "black")};
    background-color: ${(props) => (props.isActive ? "#3333cc" : props.isLoading ? "#9999ff" : "white")};

    &:hover {
        color: ${(props) => (props.isActive ? "white" : "black")};
        background-color: ${(props) => (props.isActive ? "#3333cc" : props.isLoading ? "#7777ee" : "#bbbbff")};
    }
`;
