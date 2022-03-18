import { Router } from "express";

export default interface RequestHandler {

    route(router: Router): void;

}