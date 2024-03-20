import { getCreator, getCreators, getProducts, getSortedCreators } from "@/grpc";
import { GetCreatorsRequest, GetProductsRequest, GetSortedCreatorsRequest, SortOrder } from '@/grpc/marketplace_pb';

export default async function Home() {
  // const request = new GetCreatorsRequest();
  // const creators = await getCreators(request);

  // const request1 = new GetProductsRequest();
  // const products = await getProducts(request1);

  const request3 = new GetSortedCreatorsRequest();
  request3.setLimit(3);
  request3.setOrder(SortOrder.ASCENDING);
  const sortedCreators = await getSortedCreators(request3);
  

  return (
    <main className="flex min-h-screen flex-col items-center justify-between p-24">
      {
        sortedCreators.getCreatorsList().map(item => item.getEmail()).map(email => (<div key={email}>{email}</div>))
      }
    </main>
  );
}
