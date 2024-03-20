import { CreatorsServiceClient, ProductsServiceClient } from './marketplace_grpc_pb';
import * as grpc from '@grpc/grpc-js';
import util from 'util';
import { 
    GetCreatorRequest,
    GetCreatorResponse,
    GetCreatorsRequest,
    GetCreatorsResponse,
    GetProductRequest,
    GetProductResponse,
    GetProductsRequest,
    GetProductsResponse,
    GetSortedCreatorsRequest,
    GetSortedCreatorsResponse
} from './marketplace_pb';

const clientCreators = new CreatorsServiceClient('localhost:50051', grpc.credentials.createInsecure());
const clientProducts = new ProductsServiceClient('localhost:50051', grpc.credentials.createInsecure());

const getCreatorPromise = util.promisify<GetCreatorRequest, GetCreatorResponse>(clientCreators.getCreator);
export const getCreator = getCreatorPromise.bind(clientCreators);

const getCreatorsPromise = util.promisify<GetCreatorsRequest, GetCreatorsResponse>(clientCreators.getCreators);
export const getCreators = getCreatorsPromise.bind(clientCreators);

const getSortedCreatorsPromise = util.promisify<GetSortedCreatorsRequest, GetSortedCreatorsResponse>(clientCreators.getSortedCreators);
export const getSortedCreators = getSortedCreatorsPromise.bind(clientCreators);

const getProductPromise = util.promisify<GetProductRequest, GetProductResponse>(clientProducts.getProduct);
export const getProduct = getProductPromise.bind(clientProducts);

const getProductsPromise = util.promisify<GetProductsRequest, GetProductsResponse>(clientProducts.getProducts);
export const getProducts = getProductsPromise.bind(clientProducts);